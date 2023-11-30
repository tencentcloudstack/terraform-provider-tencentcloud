package tencentcloud

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCynosdbParamTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbParamTemplatesRead,
		Schema: map[string]*schema.Schema{
			"engine_versions": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Database engine version number.",
			},

			"template_names": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The name list of templates.",
			},

			"template_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "The id list of templates.",
			},

			"db_modes": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Database mode, optional values: NORMAL, SERVERLESS.",
			},

			"offset": {
				Optional:    true,
				Default:     0,
				Type:        schema.TypeInt,
				Description: "Page offset.",
			},

			"limit": {
				Optional:    true,
				Default:     10,
				Type:        schema.TypeInt,
				Description: "Query limit.",
			},

			"products": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The product type corresponding to the query template.",
			},

			"template_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Template types.",
			},

			"engine_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Engine types.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The sort field for the returned results.",
			},

			"order_direction": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort by (asc, desc).",
			},

			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Parameter Template Information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID of template.",
						},
						"template_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of template.",
						},
						"template_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of template.",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine version.",
						},
						"db_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database mode, optional values: NORMAL, SERVERLESS.",
						},
						"param_info_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Parameter template details.Note: This field may return null, indicating that no valid value can be obtained.",
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
										Description: "An optional set of value types when the parameter type is enum.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"max": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The maximum value when the parameter type is float/integer.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"min": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The minimum value when the parameter type is float/integer.Note: This field may return null, indicating that no valid value can be obtained.",
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
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The description of parameter.",
									},
									"param_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Parameter type: integer/float/string/enum.",
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

func dataSourceTencentCloudCynosdbParamTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cynosdb_param_templates.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("engine_versions"); ok {
		engineVersionsSet := v.(*schema.Set).List()
		paramMap["engine_versions"] = helper.InterfacesStringsPoint(engineVersionsSet)
	}

	if v, ok := d.GetOk("template_names"); ok {
		templateNamesSet := v.(*schema.Set).List()
		paramMap["template_names"] = helper.InterfacesStringsPoint(templateNamesSet)
	}

	if v, ok := d.GetOk("template_ids"); ok {
		templateIds := make([]*int64, 0)
		for _, item := range v.(*schema.Set).List() {
			templateIds = append(templateIds, helper.IntInt64(item.(int)))
		}
		paramMap["template_ids"] = templateIds

	}

	if v, ok := d.GetOk("db_modes"); ok {
		dbModesSet := v.(*schema.Set).List()
		paramMap["db_modes"] = helper.InterfacesStringsPoint(dbModesSet)
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["offset"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["limit"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("products"); ok {
		productsSet := v.(*schema.Set).List()
		paramMap["products"] = helper.InterfacesStringsPoint(productsSet)
	}

	if v, ok := d.GetOk("template_types"); ok {
		templateTypesSet := v.(*schema.Set).List()
		paramMap["template_types"] = helper.InterfacesStringsPoint(templateTypesSet)
	}

	if v, ok := d.GetOk("engine_types"); ok {
		engineTypesSet := v.(*schema.Set).List()
		paramMap["engine_types"] = helper.InterfacesStringsPoint(engineTypesSet)
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["order_by"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_direction"); ok {
		paramMap["order_direction"] = helper.String(v.(string))
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var items []*cynosdb.ParamTemplateListInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbParamTemplatesByFilter(ctx, paramMap)
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
		ids = append(ids, strconv.FormatInt(*item.Id, 10))
		itemMap := make(map[string]interface{})
		itemMap["id"] = item.Id
		itemMap["template_name"] = item.TemplateName
		itemMap["template_description"] = item.TemplateDescription
		itemMap["engine_version"] = item.EngineVersion
		itemMap["db_mode"] = item.DbMode
		paramInfos := make([]map[string]interface{}, 0)
		if item.ParamInfoSet != nil {
			for _, paramInfo := range item.ParamInfoSet {
				paramInfoMap := make(map[string]interface{})
				paramInfoMap["current_value"] = paramInfo.CurrentValue
				paramInfoMap["default"] = paramInfo.Default
				enumValues := make([]string, 0)
				if paramInfo.EnumValue != nil {
					for _, enumValue := range paramInfo.EnumValue {
						enumValues = append(enumValues, *enumValue)
					}
				}
				paramInfoMap["enum_value"] = enumValues
				paramInfoMap["max"] = paramInfo.Max
				paramInfoMap["min"] = paramInfo.Min
				paramInfoMap["param_name"] = paramInfo.ParamName
				paramInfoMap["need_reboot"] = paramInfo.NeedReboot
				paramInfoMap["description"] = paramInfo.Description
				paramInfoMap["param_type"] = paramInfo.ParamType

				paramInfos = append(paramInfos, paramInfoMap)
			}
		}
		itemMap["param_info_set"] = paramInfos
		itemMap["id"] = item.Id
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
