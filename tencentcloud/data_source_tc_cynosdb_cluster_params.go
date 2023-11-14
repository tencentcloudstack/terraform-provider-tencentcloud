/*
Use this data source to query detailed information of cynosdb cluster_params

Example Usage

```hcl
data "tencentcloud_cynosdb_cluster_params" "cluster_params" {
  cluster_id = &lt;nil&gt;
  param_name = &lt;nil&gt;
  total_count = &lt;nil&gt;
  items {
		current_value = &lt;nil&gt;
		default = &lt;nil&gt;
		enum_value = &lt;nil&gt;
		max = &lt;nil&gt;
		min = &lt;nil&gt;
		param_name = &lt;nil&gt;
		need_reboot = &lt;nil&gt;
		param_type = &lt;nil&gt;
		match_type = &lt;nil&gt;
		match_value = &lt;nil&gt;
		description = &lt;nil&gt;
		is_global = &lt;nil&gt;
		modifiable_info = &lt;nil&gt;
		is_func = &lt;nil&gt;
		func = &lt;nil&gt;

  }
}
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCynosdbClusterParams() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbClusterParamsRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of cluster.",
			},

			"param_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Parameter name.",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Description: "The total count of parameters.",
			},

			"items": {
				Type:        schema.TypeList,
				Description: "Instance parameter list.Note: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"current_value": {
							Type:        schema.TypeString,
							Description: "Current value.",
						},
						"default": {
							Type:        schema.TypeString,
							Description: "Default value.",
						},
						"enum_value": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "When the parameter is enum/string/bool, the optional value list.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"max": {
							Type:        schema.TypeString,
							Description: "The maximum value when the parameter type is float/integer.",
						},
						"min": {
							Type:        schema.TypeString,
							Description: "The minimum value when the parameter type is float/integer.",
						},
						"param_name": {
							Type:        schema.TypeString,
							Description: "The name of parameter.",
						},
						"need_reboot": {
							Type:        schema.TypeInt,
							Description: "Whether to reboot.",
						},
						"param_type": {
							Type:        schema.TypeString,
							Description: "Parameter type: integer/float/string/enum/bool.",
						},
						"match_type": {
							Type:        schema.TypeString,
							Description: "Matching type, multiVal, regex is used when the parameter type is string.",
						},
						"match_value": {
							Type:        schema.TypeString,
							Description: "Match the target value, when multiVal, each key is divided by &amp;#39;;&amp;#39;.",
						},
						"description": {
							Type:        schema.TypeString,
							Description: "The description of parameter.",
						},
						"is_global": {
							Type:        schema.TypeInt,
							Description: "Is it a global parameter.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"modifiable_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Whether the parameter can be modified.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"is_func": {
							Type:        schema.TypeBool,
							Description: "Is it a function.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"func": {
							Type:        schema.TypeString,
							Description: "Function.Note: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudCynosdbClusterParamsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cynosdb_cluster_params.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("param_name"); ok {
		paramMap["ParamName"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("total_count"); v != nil {
		paramMap["TotalCount"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("items"); ok {
		itemsSet := v.([]interface{})
		tmpSet := make([]*cynosdb.ParamInfo, 0, len(itemsSet))

		for _, item := range itemsSet {
			paramInfo := cynosdb.ParamInfo{}
			paramInfoMap := item.(map[string]interface{})

			if v, ok := paramInfoMap["current_value"]; ok {
				paramInfo.CurrentValue = helper.String(v.(string))
			}
			if v, ok := paramInfoMap["default"]; ok {
				paramInfo.Default = helper.String(v.(string))
			}
			if v, ok := paramInfoMap["enum_value"]; ok {
				enumValueSet := v.(*schema.Set).List()
				paramInfo.EnumValue = helper.InterfacesStringsPoint(enumValueSet)
			}
			if v, ok := paramInfoMap["max"]; ok {
				paramInfo.Max = helper.String(v.(string))
			}
			if v, ok := paramInfoMap["min"]; ok {
				paramInfo.Min = helper.String(v.(string))
			}
			if v, ok := paramInfoMap["param_name"]; ok {
				paramInfo.ParamName = helper.String(v.(string))
			}
			if v, ok := paramInfoMap["need_reboot"]; ok {
				paramInfo.NeedReboot = helper.IntInt64(v.(int))
			}
			if v, ok := paramInfoMap["param_type"]; ok {
				paramInfo.ParamType = helper.String(v.(string))
			}
			if v, ok := paramInfoMap["match_type"]; ok {
				paramInfo.MatchType = helper.String(v.(string))
			}
			if v, ok := paramInfoMap["match_value"]; ok {
				paramInfo.MatchValue = helper.String(v.(string))
			}
			if v, ok := paramInfoMap["description"]; ok {
				paramInfo.Description = helper.String(v.(string))
			}
			if v, ok := paramInfoMap["is_global"]; ok {
				paramInfo.IsGlobal = helper.IntInt64(v.(int))
			}
			if modifiableInfoMap, ok := helper.InterfaceToMap(paramInfoMap, "modifiable_info"); ok {
				modifiableInfo := cynosdb.ModifiableInfo{}
				paramInfo.ModifiableInfo = &modifiableInfo
			}
			if v, ok := paramInfoMap["is_func"]; ok {
				paramInfo.IsFunc = helper.Bool(v.(bool))
			}
			if v, ok := paramInfoMap["func"]; ok {
				paramInfo.Func = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &paramInfo)
		}
		paramMap["items"] = tmpSet
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var items []*cynosdb.ParamInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbClusterParamsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		items = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(items))
	tmpList := make([]map[string]interface{}, 0, len(items))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
