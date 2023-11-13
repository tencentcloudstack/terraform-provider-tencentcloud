/*
Use this data source to query detailed information of mongodb instance_params

Example Usage

```hcl
data "tencentcloud_mongodb_instance_params" "instance_params" {
  instance_id = "cmgo-9d0p6umb"
        }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMongodbInstanceParams() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMongodbInstanceParamsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "InstanceId.",
			},

			"instance_enum_param": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Enum parameter.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"current_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current value.",
						},
						"default_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Default value.",
						},
						"enum_value": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Enumvalue.",
						},
						"need_restart": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "If need restart.",
						},
						"param_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of parameter.",
						},
						"tips": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Descripition of parameter.",
						},
						"value_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Value type.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "If is running.",
						},
					},
				},
			},

			"instance_integer_param": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Integer parameter.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"current_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current value.",
						},
						"default_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Default value.",
						},
						"max": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Max value.",
						},
						"min": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Min value.",
						},
						"need_restart": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "If need restart.",
						},
						"param_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of parameter.",
						},
						"tips": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Descripition of parameter.",
						},
						"value_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Value type.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "If is running.",
						},
					},
				},
			},

			"instance_text_param": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Text parameter.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"current_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current value.",
						},
						"default_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Default value.",
						},
						"need_restart": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "If need restart.",
						},
						"param_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of parameter.",
						},
						"text_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Text value.",
						},
						"tips": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Descripition of parameter.",
						},
						"value_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Value type.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "If is running.",
						},
					},
				},
			},

			"instance_multi_param": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Multi parameter.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"current_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current value.",
						},
						"default_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Default value.",
						},
						"enum_value": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Enum value.",
						},
						"need_restart": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "If need restart.",
						},
						"param_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of parameter.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "If is running.",
						},
						"tips": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Descripition of parameter.",
						},
						"value_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Value type.",
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

func dataSourceTencentCloudMongodbInstanceParamsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mongodb_instance_params.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := MongodbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var instanceEnumParam []*mongodb.InstanceEnumParam

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMongodbInstanceParamsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		instanceEnumParam = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceEnumParam))
	tmpList := make([]map[string]interface{}, 0, len(instanceEnumParam))

	if instanceEnumParam != nil {
		for _, instanceEnumParam := range instanceEnumParam {
			instanceEnumParamMap := map[string]interface{}{}

			if instanceEnumParam.CurrentValue != nil {
				instanceEnumParamMap["current_value"] = instanceEnumParam.CurrentValue
			}

			if instanceEnumParam.DefaultValue != nil {
				instanceEnumParamMap["default_value"] = instanceEnumParam.DefaultValue
			}

			if instanceEnumParam.EnumValue != nil {
				instanceEnumParamMap["enum_value"] = instanceEnumParam.EnumValue
			}

			if instanceEnumParam.NeedRestart != nil {
				instanceEnumParamMap["need_restart"] = instanceEnumParam.NeedRestart
			}

			if instanceEnumParam.ParamName != nil {
				instanceEnumParamMap["param_name"] = instanceEnumParam.ParamName
			}

			if instanceEnumParam.Tips != nil {
				instanceEnumParamMap["tips"] = instanceEnumParam.Tips
			}

			if instanceEnumParam.ValueType != nil {
				instanceEnumParamMap["value_type"] = instanceEnumParam.ValueType
			}

			if instanceEnumParam.Status != nil {
				instanceEnumParamMap["status"] = instanceEnumParam.Status
			}

			ids = append(ids, *instanceEnumParam.InstanceId)
			tmpList = append(tmpList, instanceEnumParamMap)
		}

		_ = d.Set("instance_enum_param", tmpList)
	}

	if instanceIntegerParam != nil {
		for _, instanceIntegerParam := range instanceIntegerParam {
			instanceIntegerParamMap := map[string]interface{}{}

			if instanceIntegerParam.CurrentValue != nil {
				instanceIntegerParamMap["current_value"] = instanceIntegerParam.CurrentValue
			}

			if instanceIntegerParam.DefaultValue != nil {
				instanceIntegerParamMap["default_value"] = instanceIntegerParam.DefaultValue
			}

			if instanceIntegerParam.Max != nil {
				instanceIntegerParamMap["max"] = instanceIntegerParam.Max
			}

			if instanceIntegerParam.Min != nil {
				instanceIntegerParamMap["min"] = instanceIntegerParam.Min
			}

			if instanceIntegerParam.NeedRestart != nil {
				instanceIntegerParamMap["need_restart"] = instanceIntegerParam.NeedRestart
			}

			if instanceIntegerParam.ParamName != nil {
				instanceIntegerParamMap["param_name"] = instanceIntegerParam.ParamName
			}

			if instanceIntegerParam.Tips != nil {
				instanceIntegerParamMap["tips"] = instanceIntegerParam.Tips
			}

			if instanceIntegerParam.ValueType != nil {
				instanceIntegerParamMap["value_type"] = instanceIntegerParam.ValueType
			}

			if instanceIntegerParam.Status != nil {
				instanceIntegerParamMap["status"] = instanceIntegerParam.Status
			}

			ids = append(ids, *instanceIntegerParam.InstanceId)
			tmpList = append(tmpList, instanceIntegerParamMap)
		}

		_ = d.Set("instance_integer_param", tmpList)
	}

	if instanceTextParam != nil {
		for _, instanceTextParam := range instanceTextParam {
			instanceTextParamMap := map[string]interface{}{}

			if instanceTextParam.CurrentValue != nil {
				instanceTextParamMap["current_value"] = instanceTextParam.CurrentValue
			}

			if instanceTextParam.DefaultValue != nil {
				instanceTextParamMap["default_value"] = instanceTextParam.DefaultValue
			}

			if instanceTextParam.NeedRestart != nil {
				instanceTextParamMap["need_restart"] = instanceTextParam.NeedRestart
			}

			if instanceTextParam.ParamName != nil {
				instanceTextParamMap["param_name"] = instanceTextParam.ParamName
			}

			if instanceTextParam.TextValue != nil {
				instanceTextParamMap["text_value"] = instanceTextParam.TextValue
			}

			if instanceTextParam.Tips != nil {
				instanceTextParamMap["tips"] = instanceTextParam.Tips
			}

			if instanceTextParam.ValueType != nil {
				instanceTextParamMap["value_type"] = instanceTextParam.ValueType
			}

			if instanceTextParam.Status != nil {
				instanceTextParamMap["status"] = instanceTextParam.Status
			}

			ids = append(ids, *instanceTextParam.InstanceId)
			tmpList = append(tmpList, instanceTextParamMap)
		}

		_ = d.Set("instance_text_param", tmpList)
	}

	if instanceMultiParam != nil {
		for _, instanceMultiParam := range instanceMultiParam {
			instanceMultiParamMap := map[string]interface{}{}

			if instanceMultiParam.CurrentValue != nil {
				instanceMultiParamMap["current_value"] = instanceMultiParam.CurrentValue
			}

			if instanceMultiParam.DefaultValue != nil {
				instanceMultiParamMap["default_value"] = instanceMultiParam.DefaultValue
			}

			if instanceMultiParam.EnumValue != nil {
				instanceMultiParamMap["enum_value"] = instanceMultiParam.EnumValue
			}

			if instanceMultiParam.NeedRestart != nil {
				instanceMultiParamMap["need_restart"] = instanceMultiParam.NeedRestart
			}

			if instanceMultiParam.ParamName != nil {
				instanceMultiParamMap["param_name"] = instanceMultiParam.ParamName
			}

			if instanceMultiParam.Status != nil {
				instanceMultiParamMap["status"] = instanceMultiParam.Status
			}

			if instanceMultiParam.Tips != nil {
				instanceMultiParamMap["tips"] = instanceMultiParam.Tips
			}

			if instanceMultiParam.ValueType != nil {
				instanceMultiParamMap["value_type"] = instanceMultiParam.ValueType
			}

			ids = append(ids, *instanceMultiParam.InstanceId)
			tmpList = append(tmpList, instanceMultiParamMap)
		}

		_ = d.Set("instance_multi_param", tmpList)
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
