/*
Use this data source to query detailed information of scf triggers

Example Usage

```hcl
data "tencentcloud_scf_triggers" "triggers" {
  function_name = "keep-1676351130"
  namespace     = "default"
  order_by      = "add_time"
  order         = "DESC"
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

func dataSourceTencentCloudScfTriggers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudScfTriggersRead,
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Function name.",
			},

			"namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Namespace. Default value: default.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Indicates by which field to sort the returned results. Valid values: add_time, mod_time. Default value: mod_time.",
			},

			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Indicates whether the returned results are sorted in ascending or descending order. Valid values: ASC, DESC. Default value: DESC.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "* Qualifier:Function version, alias.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Fields to be filtered. Up to 10 conditions allowed.Values of Name: VpcId, SubnetId, ClsTopicId, ClsLogsetId, Role, CfsId, CfsMountInsId, Eip. Values limit: 1.Name options: Status, Runtime, FunctionType, PublicNetStatus, AsyncRunEnable, TraceEnable. Values limit: 20.When Name is Runtime, CustomImage refers to the image type function.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Filter values of the field.",
						},
					},
				},
			},

			"triggers": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Trigger list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to enable.",
						},
						"qualifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Function version or alias.",
						},
						"trigger_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Trigger name.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Trigger type.",
						},
						"trigger_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Detailed configuration of trigger.",
						},
						"available_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether the trigger is available.",
						},
						"custom_argument": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom parameterNote: this field may return null, indicating that no valid values can be obtained.",
						},
						"add_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Trigger creation time.",
						},
						"mod_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Trigger last modified time.",
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Minimum resource ID of trigger.",
						},
						"bind_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Trigger-Function binding status.",
						},
						"trigger_attribute": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Trigger type. Two-way means that the trigger can be manipulated in both consoles, while one-way means that the trigger can be created only in the SCF Console.",
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

func dataSourceTencentCloudScfTriggersRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_scf_triggers.read")()
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

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order"); ok {
		paramMap["Order"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*scf.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := scf.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["Filters"] = tmpSet
	}

	service := ScfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var triggers []*scf.TriggerInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeScfTriggersByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		triggers = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(triggers))
	tmpList := make([]map[string]interface{}, 0, len(triggers))

	if triggers != nil {
		for _, triggerInfo := range triggers {
			triggerInfoMap := map[string]interface{}{}

			if triggerInfo.Enable != nil {
				triggerInfoMap["enable"] = triggerInfo.Enable
			}

			if triggerInfo.Qualifier != nil {
				triggerInfoMap["qualifier"] = triggerInfo.Qualifier
			}

			if triggerInfo.TriggerName != nil {
				triggerInfoMap["trigger_name"] = triggerInfo.TriggerName
			}

			if triggerInfo.Type != nil {
				triggerInfoMap["type"] = triggerInfo.Type
			}

			if triggerInfo.TriggerDesc != nil {
				triggerInfoMap["trigger_desc"] = triggerInfo.TriggerDesc
			}

			if triggerInfo.AvailableStatus != nil {
				triggerInfoMap["available_status"] = triggerInfo.AvailableStatus
			}

			if triggerInfo.CustomArgument != nil {
				triggerInfoMap["custom_argument"] = triggerInfo.CustomArgument
			}

			if triggerInfo.AddTime != nil {
				triggerInfoMap["add_time"] = triggerInfo.AddTime
			}

			if triggerInfo.ModTime != nil {
				triggerInfoMap["mod_time"] = triggerInfo.ModTime
			}

			if triggerInfo.ResourceId != nil {
				triggerInfoMap["resource_id"] = triggerInfo.ResourceId
			}

			if triggerInfo.BindStatus != nil {
				triggerInfoMap["bind_status"] = triggerInfo.BindStatus
			}

			if triggerInfo.TriggerAttribute != nil {
				triggerInfoMap["trigger_attribute"] = triggerInfo.TriggerAttribute
			}

			ids = append(ids, *triggerInfo.TriggerName)
			tmpList = append(tmpList, triggerInfoMap)
		}

		_ = d.Set("triggers", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
