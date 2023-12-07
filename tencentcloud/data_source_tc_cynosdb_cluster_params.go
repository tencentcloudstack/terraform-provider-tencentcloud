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

			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Instance parameter list. Note: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"current_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current value.",
						},
						"default": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Default value.",
						},
						"enum_value": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "When the parameter is enum/string/bool, the optional value list.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"max": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The maximum value when the parameter type is float/integer.",
						},
						"min": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The minimum value when the parameter type is float/integer.",
						},
						"param_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of parameter.",
						},
						"need_reboot": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to reboot.",
						},
						"param_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter type: integer/float/string/enum/bool.",
						},
						"match_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Matching type, multiVal, regex is used when the parameter type is string.",
						},
						"match_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Match the target value, when multiVal, each key is divided by `;`.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of parameter.",
						},
						"is_global": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Is it a global parameter.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"is_func": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Is it a function.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"func": {
							Type:        schema.TypeString,
							Computed:    true,
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

	var clusterId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("param_name"); ok {
		paramMap["param_name"] = v.(string)
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var items []*cynosdb.ParamInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClusterParamsByFilter(ctx, clusterId, paramMap)
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
	for _, item := range items {
		ids = append(ids, *item.ParamName)
		itemMap := make(map[string]interface{})
		itemMap["current_value"] = item.CurrentValue
		itemMap["default"] = item.Default
		itemMap["max"] = item.Max
		itemMap["min"] = item.Min
		itemMap["param_name"] = item.ParamName
		itemMap["need_reboot"] = item.NeedReboot
		itemMap["param_type"] = item.ParamType
		itemMap["match_type"] = item.MatchType
		itemMap["match_value"] = item.MatchValue
		itemMap["description"] = item.Description
		itemMap["is_global"] = item.IsGlobal
		itemMap["is_func"] = item.IsFunc
		itemMap["func"] = item.Func
		enumValues := make([]string, 0)
		if item.EnumValue != nil {
			for _, enumValueItem := range item.EnumValue {
				enumValues = append(enumValues, *enumValueItem)
			}
			itemMap["enum_value"] = enumValues
		}
		tmpList = append(tmpList, itemMap)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("items", tmpList)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
