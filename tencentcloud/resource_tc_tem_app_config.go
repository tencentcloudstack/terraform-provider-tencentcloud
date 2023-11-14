/*
Provides a resource to create a tem app_config

Example Usage

```hcl
resource "tencentcloud_tem_app_config" "app_config" {
  environment_id = "en-xxx"
  name = "xxx"
  data {
		key = "key"
		value = "value"

  }
}
```

Import

tem app_config can be imported using the id, e.g.

```
terraform import tencentcloud_tem_app_config.app_config app_config_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tem "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tem/v20210701"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudTemAppConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTemAppConfigCreate,
		Read:   resourceTencentCloudTemAppConfigRead,
		Update: resourceTencentCloudTemAppConfigUpdate,
		Delete: resourceTencentCloudTemAppConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"environment_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Environment ID.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "AppConfig name.",
			},

			"data": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Payload.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTemAppConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_app_config.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = tem.NewCreateConfigDataRequest()
		response      = tem.NewCreateConfigDataResponse()
		environmentId string
		name          string
	)
	if v, ok := d.GetOk("environment_id"); ok {
		environmentId = v.(string)
		request.EnvironmentId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("data"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			pair := tem.Pair{}
			if v, ok := dMap["key"]; ok {
				pair.Key = helper.String(v.(string))
			}
			if v, ok := dMap["value"]; ok {
				pair.Value = helper.String(v.(string))
			}
			request.Data = append(request.Data, &pair)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTemClient().CreateConfigData(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tem appConfig failed, reason:%+v", logId, err)
		return err
	}

	environmentId = *response.Response.EnvironmentId
	d.SetId(strings.Join([]string{environmentId, name}, FILED_SP))

	return resourceTencentCloudTemAppConfigRead(d, meta)
}

func resourceTencentCloudTemAppConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_app_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TemService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	name := idSplit[1]

	appConfig, err := service.DescribeTemAppConfigById(ctx, environmentId, name)
	if err != nil {
		return err
	}

	if appConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TemAppConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if appConfig.EnvironmentId != nil {
		_ = d.Set("environment_id", appConfig.EnvironmentId)
	}

	if appConfig.Name != nil {
		_ = d.Set("name", appConfig.Name)
	}

	if appConfig.Data != nil {
		dataList := []interface{}{}
		for _, data := range appConfig.Data {
			dataMap := map[string]interface{}{}

			if appConfig.Data.Key != nil {
				dataMap["key"] = appConfig.Data.Key
			}

			if appConfig.Data.Value != nil {
				dataMap["value"] = appConfig.Data.Value
			}

			dataList = append(dataList, dataMap)
		}

		_ = d.Set("data", dataList)

	}

	return nil
}

func resourceTencentCloudTemAppConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_app_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tem.NewModifyConfigDataRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	name := idSplit[1]

	request.EnvironmentId = &environmentId
	request.Name = &name

	immutableArgs := []string{"environment_id", "name", "data"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("data") {
		if v, ok := d.GetOk("data"); ok {
			for _, item := range v.([]interface{}) {
				pair := tem.Pair{}
				if v, ok := dMap["key"]; ok {
					pair.Key = helper.String(v.(string))
				}
				if v, ok := dMap["value"]; ok {
					pair.Value = helper.String(v.(string))
				}
				request.Data = append(request.Data, &pair)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTemClient().ModifyConfigData(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tem appConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTemAppConfigRead(d, meta)
}

func resourceTencentCloudTemAppConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_app_config.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TemService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	name := idSplit[1]

	if err := service.DeleteTemAppConfigById(ctx, environmentId, name); err != nil {
		return err
	}

	return nil
}
