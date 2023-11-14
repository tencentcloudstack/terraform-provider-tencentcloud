/*
Use this data source to query detailed information of scf provisioned_concurrency_config

Example Usage

```hcl
data "tencentcloud_scf_provisioned_concurrency_config" "provisioned_concurrency_config" {
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
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudScfProvisionedConcurrencyConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudScfProvisionedConcurrencyConfigRead,
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Name of the function for which to get the provisioned concurrency details.",
			},

			"namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Function namespace. Default value: default.",
			},

			"qualifier": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Function version number. If this parameter is left empty, the provisioned concurrency information of all function versions will be returned.",
			},

			"unallocated_concurrency_num": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Unallocated provisioned concurrency amount of function.",
			},

			"allocated": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Allocated provisioned concurrency amount of function.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allocated_provisioned_concurrency_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Set provisioned concurrency amount.",
						},
						"available_provisioned_concurrency_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Currently available provisioned concurrency amount.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Provisioned concurrency setting task status. `Done`: completed; `InProgress`: in progress; `Failed`: partially or completely failed.",
						},
						"status_reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status description of provisioned concurrency setting task.",
						},
						"qualifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Function version number.",
						},
						"trigger_actions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of scheduled provisioned concurrency scaling actionsNote: this field may return `null`, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"trigger_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Scheduled action nameNote: this field may return `null`, indicating that no valid values can be obtained.",
									},
									"trigger_provisioned_concurrency_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Target provisioned concurrency of the scheduled scaling action Note: this field may return `null`, indicating that no valid values can be obtained.",
									},
									"trigger_cron_config": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Trigger time of the scheduled action in Cron expression. Seven fields are required and should be separated with a space.Note: this field may return `null`, indicating that no valid values can be obtained.",
									},
									"provisioned_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The provision type. Value: `Default`Note: This field may return `null`, indicating that no valid value can be found.",
									},
								},
							},
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

func dataSourceTencentCloudScfProvisionedConcurrencyConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_scf_provisioned_concurrency_config.read")()
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

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeScfProvisionedConcurrencyConfigByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		unallocatedConcurrencyNum = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(unallocatedConcurrencyNum))
	if unallocatedConcurrencyNum != nil {
		_ = d.Set("unallocated_concurrency_num", unallocatedConcurrencyNum)
	}

	if allocated != nil {
		for _, versionProvisionedConcurrencyInfo := range allocated {
			versionProvisionedConcurrencyInfoMap := map[string]interface{}{}

			if versionProvisionedConcurrencyInfo.AllocatedProvisionedConcurrencyNum != nil {
				versionProvisionedConcurrencyInfoMap["allocated_provisioned_concurrency_num"] = versionProvisionedConcurrencyInfo.AllocatedProvisionedConcurrencyNum
			}

			if versionProvisionedConcurrencyInfo.AvailableProvisionedConcurrencyNum != nil {
				versionProvisionedConcurrencyInfoMap["available_provisioned_concurrency_num"] = versionProvisionedConcurrencyInfo.AvailableProvisionedConcurrencyNum
			}

			if versionProvisionedConcurrencyInfo.Status != nil {
				versionProvisionedConcurrencyInfoMap["status"] = versionProvisionedConcurrencyInfo.Status
			}

			if versionProvisionedConcurrencyInfo.StatusReason != nil {
				versionProvisionedConcurrencyInfoMap["status_reason"] = versionProvisionedConcurrencyInfo.StatusReason
			}

			if versionProvisionedConcurrencyInfo.Qualifier != nil {
				versionProvisionedConcurrencyInfoMap["qualifier"] = versionProvisionedConcurrencyInfo.Qualifier
			}

			if versionProvisionedConcurrencyInfo.TriggerActions != nil {
				triggerActionsList := []interface{}{}
				for _, triggerActions := range versionProvisionedConcurrencyInfo.TriggerActions {
					triggerActionsMap := map[string]interface{}{}

					if triggerActions.TriggerName != nil {
						triggerActionsMap["trigger_name"] = triggerActions.TriggerName
					}

					if triggerActions.TriggerProvisionedConcurrencyNum != nil {
						triggerActionsMap["trigger_provisioned_concurrency_num"] = triggerActions.TriggerProvisionedConcurrencyNum
					}

					if triggerActions.TriggerCronConfig != nil {
						triggerActionsMap["trigger_cron_config"] = triggerActions.TriggerCronConfig
					}

					if triggerActions.ProvisionedType != nil {
						triggerActionsMap["provisioned_type"] = triggerActions.ProvisionedType
					}

					triggerActionsList = append(triggerActionsList, triggerActionsMap)
				}

				versionProvisionedConcurrencyInfoMap["trigger_actions"] = []interface{}{triggerActionsList}
			}

			ids = append(ids, *versionProvisionedConcurrencyInfo.FunctionName)
			tmpList = append(tmpList, versionProvisionedConcurrencyInfoMap)
		}

		_ = d.Set("allocated", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
