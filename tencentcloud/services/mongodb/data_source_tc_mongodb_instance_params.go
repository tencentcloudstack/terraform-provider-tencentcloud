package mongodb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMongodbInstanceParams() *schema.Resource {
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
							Description: "current value.",
						},
						"default_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "default value.",
						},
						"enum_value": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "enumvalue.",
						},
						"need_restart": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "if need restart.",
						},
						"param_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of parameter.",
						},
						"tips": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "descripition of parameter.",
						},
						"value_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "value type.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "if is running.",
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
							Description: "current value.",
						},
						"default_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "default value.",
						},
						"max": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "max value.",
						},
						"min": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "min value.",
						},
						"need_restart": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "if need restart.",
						},
						"param_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of parameter.",
						},
						"tips": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "descripition of parameter.",
						},
						"value_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "value type.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "if is running.",
						},
					},
				},
			},

			"instance_text_param": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "text parameter.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"current_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "current value.",
						},
						"default_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "default value.",
						},
						"need_restart": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "if need restart.",
						},
						"param_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of parameter.",
						},
						"text_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "text value.",
						},
						"tips": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "descripition of parameter.",
						},
						"value_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "value type.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "if is running.",
						},
					},
				},
			},

			"instance_multi_param": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "multi parameter.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"current_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "current value.",
						},
						"default_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "default value.",
						},
						"enum_value": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "enum value.",
						},
						"need_restart": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "if need restart.",
						},
						"param_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of parameter.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "if is running.",
						},
						"tips": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "descripition of parameter.",
						},
						"value_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "value type.",
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
	defer tccommon.LogElapsed("data_source.tencentcloud_mongodb_instance_params.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	service := MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var instanceParam *mongodb.DescribeInstanceParamsResponseParams

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMongodbInstanceParams(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instanceParam = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0)

	paramList := make([]map[string]interface{}, 0)

	if instanceParam != nil {
		enumTmpList := make([]map[string]interface{}, 0, len(instanceParam.InstanceEnumParam))
		for _, instanceEnumParam := range instanceParam.InstanceEnumParam {
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

			ids = append(ids, *instanceEnumParam.ParamName)
			enumTmpList = append(enumTmpList, instanceEnumParamMap)
		}

		paramList = append(paramList, enumTmpList...)
		_ = d.Set("instance_enum_param", enumTmpList)
	}

	if instanceParam != nil {
		integerTmpList := make([]map[string]interface{}, 0, len(instanceParam.InstanceIntegerParam))
		for _, instanceIntegerParam := range instanceParam.InstanceIntegerParam {
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

			ids = append(ids, *instanceIntegerParam.ParamName)
			integerTmpList = append(integerTmpList, instanceIntegerParamMap)
		}

		paramList = append(paramList, integerTmpList...)

		_ = d.Set("instance_integer_param", integerTmpList)
	}

	if instanceParam != nil {
		enumTextList := make([]map[string]interface{}, 0, len(instanceParam.InstanceTextParam))
		for _, instanceTextParam := range instanceParam.InstanceTextParam {
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

			ids = append(ids, *instanceTextParam.ParamName)
			enumTextList = append(enumTextList, instanceTextParamMap)
		}

		paramList = append(paramList, enumTextList...)
		_ = d.Set("instance_text_param", enumTextList)
	}

	if instanceParam != nil {
		enumMultiList := make([]map[string]interface{}, 0, len(instanceParam.InstanceMultiParam))
		for _, instanceMultiParam := range instanceParam.InstanceMultiParam {
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

			ids = append(ids, *instanceMultiParam.ParamName)
			enumMultiList = append(enumMultiList, instanceMultiParamMap)
		}
		paramList = append(paramList, enumMultiList...)
		_ = d.Set("instance_multi_param", enumMultiList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), paramList); e != nil {
			return e
		}
	}
	return nil
}
