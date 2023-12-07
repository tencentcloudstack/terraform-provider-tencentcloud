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

func resourceTencentCloudClickhouseXmlConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClickhouseXmlConfigCreate,
		Read:   resourceTencentCloudClickhouseXmlConfigRead,
		Update: resourceTencentCloudClickhouseXmlConfigUpdate,
		Delete: resourceTencentCloudClickhouseXmlConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"modify_conf_context": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Configuration file modification information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration file name.",
						},
						"new_conf_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "New content of configuration file, base64 encoded.",
						},
						"file_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Path to save configuration file.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudClickhouseXmlConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clickhouse_xml_config.create")()
	defer inconsistentCheck(d, meta)()

	var ids []string
	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}
	ids = append(ids, instanceId)

	if row, ok := d.GetOk("modify_conf_context"); ok {
		items := row.([]interface{})
		for _, v := range items {
			value := v.(map[string]interface{})
			fileName := value["file_name"].(string)

			ids = append(ids, fileName)
		}
	}

	d.SetId(strings.Join(ids, FILED_SP))

	return resourceTencentCloudClickhouseXmlConfigUpdate(d, meta)
}

func resourceTencentCloudClickhouseXmlConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clickhouse_xml_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdwchService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	fileName := idSplit[1]

	xmlConfig, err := service.DescribeClickhouseXmlConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if xmlConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClickhouseXmlConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	var modifyConfContextList []interface{}
	for _, modifyConfContext := range xmlConfig {
		if fileName == *modifyConfContext.FileName {
			modifyConfContextMap := map[string]interface{}{}

			if modifyConfContext.FileName != nil {
				modifyConfContextMap["file_name"] = modifyConfContext.FileName
			}

			if modifyConfContext.OriParam != nil {
				modifyConfContextMap["new_conf_value"] = modifyConfContext.OriParam
			}

			if modifyConfContext.FilePath != nil {
				modifyConfContextMap["file_path"] = modifyConfContext.FilePath
			}

			modifyConfContextList = append(modifyConfContextList, modifyConfContextMap)
		}
	}
	_ = d.Set("modify_conf_context", modifyConfContextList)

	return nil
}

func resourceTencentCloudClickhouseXmlConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clickhouse_xml_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cdwch.NewModifyClusterConfigsRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]

	request.InstanceId = &instanceId

	var modifyConfContexts []*cdwch.ConfigSubmitContext
	if d.HasChange("modify_conf_context") {
		if row, ok := d.GetOk("modify_conf_context"); ok {
			configContexts := row.([]interface{})
			for _, v := range configContexts {
				value := v.(map[string]interface{})
				fileName := value["file_name"].(string)
				newConfValue := value["new_conf_value"].(string)
				filePath := value["file_path"].(string)

				modifyConfContexts = append(modifyConfContexts, &cdwch.ConfigSubmitContext{
					FileName:     &fileName,
					NewConfValue: &newConfValue,
					FilePath:     &filePath,
				})
			}
		}
	}
	request.ModifyConfContext = modifyConfContexts

	if len(modifyConfContexts) > 0 {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdwchClient().ModifyClusterConfigs(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cdwch xmlConfig failed, reason:%+v", logId, err)
			return err
		}
	}

	service := CdwchService{client: meta.(*TencentCloudClient).apiV3Conn}
	conf := BuildStateChangeConf([]string{}, []string{"Serving"}, 10*readRetryTimeout, time.Second, service.InstanceStateRefreshFunc(instanceId))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudClickhouseXmlConfigRead(d, meta)
}

func resourceTencentCloudClickhouseXmlConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clickhouse_xml_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
