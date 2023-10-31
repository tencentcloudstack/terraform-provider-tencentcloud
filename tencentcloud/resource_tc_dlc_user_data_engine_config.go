/*
Provides a resource to create a dlc user_data_engine_config

Example Usage

```hcl
resource "tencentcloud_dlc_user_data_engine_config" "user_data_engine_config" {
  data_engine_id = "DataEngine-cgkvbas6"
  data_engine_config_pairs {
		config_item = "qq"
		config_value = "ff"
  }
  session_resource_template {
		driver_size = "small"
		executor_size = "small"
		executor_nums = 1
		executor_max_numbers = 1
  }
}
```

Import

dlc user_data_engine_config can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_user_data_engine_config.user_data_engine_config user_data_engine_config_id
```
*/
package tencentcloud

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDlcUserDataEngineConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcUserDataEngineConfigCreate,
		Read:   resourceTencentCloudDlcUserDataEngineConfigRead,
		Update: resourceTencentCloudDlcUserDataEngineConfigUpdate,
		Delete: resourceTencentCloudDlcUserDataEngineConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"data_engine_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Engine unique id.",
			},

			"data_engine_config_pairs": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Engine configuration items.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_item": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Config key.",
						},
						"config_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Config value.",
						},
					},
				},
			},

			"session_resource_template": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Job engine resource configuration template.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"driver_size": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Engine driver size specification only supports: small/medium/large/xlarge/m.small/m.medium/m.large/m.xlarge.",
						},
						"executor_size": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Engine executor size specification only supports: small/medium/large/xlarge/m.small/m.medium/m.large/m.xlarge.",
						},
						"executor_nums": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specify the number of executors. The minimum value is 1 and the maximum value is less than the cluster specification.",
						},
						"executor_max_numbers": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specify the executor max number (in a dynamic configuration scenario), the minimum value is 1, and the maximum value is less than the cluster specification (when ExecutorMaxNumbers is less than ExecutorNums, the value is set to ExecutorNums).",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudDlcUserDataEngineConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_user_data_engine_config.create")()
	defer inconsistentCheck(d, meta)()

	return resourceTencentCloudDlcUserDataEngineConfigUpdate(d, meta)
}

func resourceTencentCloudDlcUserDataEngineConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_user_data_engine_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DlcService{client: meta.(*TencentCloudClient).apiV3Conn}

	dataEngineId := d.Id()

	userDataEngineConfig, err := service.DescribeDlcUserDataEngineConfigById(ctx, dataEngineId)
	if err != nil {
		return err
	}

	if userDataEngineConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DlcUserDataEngineConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if userDataEngineConfig.DataEngineId != nil {
		_ = d.Set("data_engine_id", userDataEngineConfig.DataEngineId)
	}

	if userDataEngineConfig.DataEngineConfigPairs != nil {
		var dataEngineConfigPairsList []interface{}
		for _, dataEngineConfigPairs := range userDataEngineConfig.DataEngineConfigPairs {
			dataEngineConfigPairsMap := map[string]interface{}{}

			if dataEngineConfigPairs.ConfigItem != nil {
				dataEngineConfigPairsMap["config_item"] = dataEngineConfigPairs.ConfigItem
			}

			if dataEngineConfigPairs.ConfigValue != nil {
				dataEngineConfigPairsMap["config_value"] = dataEngineConfigPairs.ConfigValue
			}

			dataEngineConfigPairsList = append(dataEngineConfigPairsList, dataEngineConfigPairsMap)
		}

		_ = d.Set("data_engine_config_pairs", dataEngineConfigPairsList)

	}

	if userDataEngineConfig.SessionResourceTemplate != nil {
		sessionResourceTemplateMap := map[string]interface{}{}

		if userDataEngineConfig.SessionResourceTemplate.DriverSize != nil {
			sessionResourceTemplateMap["driver_size"] = userDataEngineConfig.SessionResourceTemplate.DriverSize
		}

		if userDataEngineConfig.SessionResourceTemplate.ExecutorSize != nil {
			sessionResourceTemplateMap["executor_size"] = userDataEngineConfig.SessionResourceTemplate.ExecutorSize
		}

		if userDataEngineConfig.SessionResourceTemplate.ExecutorNums != nil {
			sessionResourceTemplateMap["executor_nums"] = userDataEngineConfig.SessionResourceTemplate.ExecutorNums
		}

		if userDataEngineConfig.SessionResourceTemplate.ExecutorMaxNumbers != nil {
			sessionResourceTemplateMap["executor_max_numbers"] = userDataEngineConfig.SessionResourceTemplate.ExecutorMaxNumbers
		}

		_ = d.Set("session_resource_template", []interface{}{sessionResourceTemplateMap})
	}

	return nil
}

func resourceTencentCloudDlcUserDataEngineConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_user_data_engine_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dlc.NewUpdateUserDataEngineConfigRequest()

	var dataEngineId string
	if v, ok := d.GetOk("data_engine_id"); ok {
		dataEngineId = v.(string)
		request.DataEngineId = helper.String(v.(string))
	}

	if d.HasChange("data_engine_config_pairs") {
		if v, ok := d.GetOk("data_engine_config_pairs"); ok {
			for _, item := range v.([]interface{}) {
				dataEngineConfigPair := dlc.DataEngineConfigPair{}
				dMap := item.(map[string]interface{})
				if v, ok := dMap["config_item"]; ok {
					dataEngineConfigPair.ConfigItem = helper.String(v.(string))
				}
				if v, ok := dMap["config_value"]; ok {
					dataEngineConfigPair.ConfigValue = helper.String(v.(string))
				}
				request.DataEngineConfigPairs = append(request.DataEngineConfigPairs, &dataEngineConfigPair)
			}
		}
	}

	if d.HasChange("session_resource_template") {
		if dMap, ok := helper.InterfacesHeadMap(d, "session_resource_template"); ok {
			sessionResourceTemplate := dlc.SessionResourceTemplate{}
			if v, ok := dMap["driver_size"]; ok {
				sessionResourceTemplate.DriverSize = helper.String(v.(string))
			}
			if v, ok := dMap["executor_size"]; ok {
				sessionResourceTemplate.ExecutorSize = helper.String(v.(string))
			}
			if v, ok := dMap["executor_nums"]; ok {
				sessionResourceTemplate.ExecutorNums = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["executor_max_numbers"]; ok {
				sessionResourceTemplate.ExecutorMaxNumbers = helper.IntUint64(v.(int))
			}
			request.SessionResourceTemplate = &sessionResourceTemplate
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDlcClient().UpdateUserDataEngineConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dlc userDataEngineConfig failed, reason:%+v", logId, err)
		return err
	}

	service := DlcService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"2"}, 5*readRetryTimeout, time.Second, service.DlcRestartDataEngineStateRefreshFunc(dataEngineId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	d.SetId(dataEngineId)
	return resourceTencentCloudDlcUserDataEngineConfigRead(d, meta)
}

func resourceTencentCloudDlcUserDataEngineConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_user_data_engine_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
