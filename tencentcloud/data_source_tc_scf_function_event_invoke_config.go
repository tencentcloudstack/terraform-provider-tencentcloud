/*
Use this data source to query detailed information of scf function_event_invoke_config

Example Usage

```hcl
data "tencentcloud_scf_function_event_invoke_config" "function_event_invoke_config" {
  function_name = ""
  namespace = ""
  qualifier = ""
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudScfFunctionEventInvokeConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudScfFunctionEventInvokeConfigRead,
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Function name.",
			},

			"namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Function namespace. Default value: default.",
			},

			"qualifier": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Function version. Default value: $LATEST.",
			},

			"async_trigger_config": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Async retry configuration information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"retry_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Async retry configuration of function upon user error.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"retry_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of retry attempts.",
									},
								},
							},
						},
						"msg_t_t_l": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Message retention period.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudScfFunctionEventInvokeConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_scf_function_event_invoke_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("function_name"); ok {
		paramMap["FunctionName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		paramMap["Namespace"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("qualifier"); ok {
		paramMap["Qualifier"] = helper.String(v.(string))
	}

	service := ScfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var asyncTriggerConfig []*scf.AsyncTriggerConfig

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeScfFunctionEventInvokeConfigByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		asyncTriggerConfig = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(asyncTriggerConfig))
	if asyncTriggerConfig != nil {
		asyncTriggerConfigMap := map[string]interface{}{}

		if asyncTriggerConfig.RetryConfig != nil {
			retryConfigList := []interface{}{}
			for _, retryConfig := range asyncTriggerConfig.RetryConfig {
				retryConfigMap := map[string]interface{}{}

				if retryConfig.RetryNum != nil {
					retryConfigMap["retry_num"] = retryConfig.RetryNum
				}

				retryConfigList = append(retryConfigList, retryConfigMap)
			}

			asyncTriggerConfigMap["retry_config"] = []interface{}{retryConfigList}
		}

		if asyncTriggerConfig.MsgTTL != nil {
			asyncTriggerConfigMap["msg_t_t_l"] = asyncTriggerConfig.MsgTTL
		}

		ids = append(ids, *asyncTriggerConfig.FunctionName)
		_ = d.Set("async_trigger_config", asyncTriggerConfigMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), asyncTriggerConfigMap); e != nil {
			return e
		}
	}
	return nil
}
