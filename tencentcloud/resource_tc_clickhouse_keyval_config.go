package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdwch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch/v20200915"
)

func resourceTencentCloudClickhouseKeyvalConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClickhouseKeyvalConfigCreate,
		Read:   resourceTencentCloudClickhouseKeyvalConfigRead,
		Update: resourceTencentCloudClickhouseKeyvalConfigUpdate,
		Delete: resourceTencentCloudClickhouseKeyvalConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"items": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "configuration list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"conf_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance config key.",
						},
						"conf_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance config value.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudClickhouseKeyvalConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clickhouse_keyval_config.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cdwch.NewModifyInstanceKeyValConfigsRequest()

	var ids []string
	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}
	ids = append(ids, instanceId)

	var addItems []*cdwch.InstanceConfigItem
	if row, ok := d.GetOk("items"); ok {
		items := row.([]interface{})
		for _, v := range items {
			value := v.(map[string]interface{})
			configKey := value["conf_key"].(string)
			configValue := value["conf_value"].(string)

			addItems = append(addItems, &cdwch.InstanceConfigItem{
				ConfKey:   &configKey,
				ConfValue: &configValue,
			})
			ids = append(ids, configKey, configValue)
		}
	}

	request.InstanceId = &instanceId
	request.AddItems = addItems

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdwchClient().ModifyInstanceKeyValConfigs(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cdwch config failed, reason:%+v", logId, err)
		return err
	}

	service := CdwchService{client: meta.(*TencentCloudClient).apiV3Conn}
	conf := BuildStateChangeConf([]string{}, []string{"Serving"}, 10*readRetryTimeout, time.Second, service.InstanceStateRefreshFunc(instanceId))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	d.SetId(strings.Join(ids, FILED_SP))

	return resourceTencentCloudClickhouseKeyvalConfigRead(d, meta)
}

func resourceTencentCloudClickhouseKeyvalConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clickhouse_keyval_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdwchService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	configKey := idSplit[1]

	_ = d.Set("instance_id", instanceId)

	configItems, err := service.DescribeClickhouseKeyvalConfigById(ctx, instanceId)
	if err != nil {
		return err
	}
	if configItems == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClickhouseKeyvalConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	resultMap := make(map[string]*cdwch.InstanceConfigInfo)
	for _, item := range configItems {
		resultMap[*item.ConfKey] = item
	}

	var itemsList []interface{}
	item := resultMap[configKey]
	if item != nil {
		itemsMap := map[string]interface{}{}
		itemsMap["conf_key"] = item.ConfKey
		itemsMap["conf_value"] = item.ConfValue
		itemsList = append(itemsList, itemsMap)
	}
	_ = d.Set("items", itemsList)

	return nil
}

func resourceTencentCloudClickhouseKeyvalConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clickhouse_keyval_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cdwch.NewModifyInstanceKeyValConfigsRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	configKey := idSplit[1]

	var ids []string
	ids = append(ids, instanceId)

	var updateItems []*cdwch.InstanceConfigItem
	if d.HasChange("items") {
		items := d.Get("items").([]interface{})
		for _, v := range items {
			value := v.(map[string]interface{})
			newConfigKey := value["conf_key"].(string)
			newConfigValue := value["conf_value"].(string)

			if configKey != newConfigKey {
				return fmt.Errorf("`conf_key` is not allowed to be modified when updating the configuration list")
			}
			updateItems = append(updateItems, &cdwch.InstanceConfigItem{
				ConfKey:   &newConfigKey,
				ConfValue: &newConfigValue,
			})
			ids = append(ids, newConfigKey, newConfigValue)
		}
	}

	request.InstanceId = &instanceId
	request.UpdateItems = updateItems

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdwchClient().ModifyInstanceKeyValConfigs(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cdwch config failed, reason:%+v", logId, err)
		return err
	}

	service := CdwchService{client: meta.(*TencentCloudClient).apiV3Conn}
	conf := BuildStateChangeConf([]string{}, []string{"Serving"}, 10*readRetryTimeout, time.Second, service.InstanceStateRefreshFunc(instanceId))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	d.SetId(strings.Join(ids, FILED_SP))

	return resourceTencentCloudClickhouseKeyvalConfigRead(d, meta)
}

func resourceTencentCloudClickhouseKeyvalConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clickhouse_keyval_config.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cdwch.NewModifyInstanceKeyValConfigsRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	configKey := idSplit[1]
	configValue := idSplit[2]

	var delItems []*cdwch.InstanceConfigItem
	delItems = append(delItems, &cdwch.InstanceConfigItem{
		ConfKey:   &configKey,
		ConfValue: &configValue,
	})

	request.InstanceId = &instanceId
	request.DelItems = delItems

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdwchClient().ModifyInstanceKeyValConfigs(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cdwch config failed, reason:%+v", logId, err)
		return err
	}

	service := CdwchService{client: meta.(*TencentCloudClient).apiV3Conn}
	conf := BuildStateChangeConf([]string{}, []string{"Serving"}, 10*readRetryTimeout, time.Second, service.InstanceStateRefreshFunc(instanceId))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
