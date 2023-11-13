/*
Use this data source to query detailed information of cynosdb param_templates

Example Usage

```hcl
data "tencentcloud_cynosdb_param_templates" "param_templates" {
  engine_versions = &lt;nil&gt;
  template_names = &lt;nil&gt;
  template_ids = &lt;nil&gt;
  db_modes = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  products = &lt;nil&gt;
  template_types = &lt;nil&gt;
  engine_types = &lt;nil&gt;
  order_by = &lt;nil&gt;
  order_direction = &lt;nil&gt;
  total_count = &lt;nil&gt;
  items {
		id = &lt;nil&gt;
		template_name = &lt;nil&gt;
		template_description = &lt;nil&gt;
		engine_version = &lt;nil&gt;
		db_mode = &lt;nil&gt;
		param_info_set {
			current_value = &lt;nil&gt;
			default = &lt;nil&gt;
			enum_value = &lt;nil&gt;
			max = &lt;nil&gt;
			min = &lt;nil&gt;
			param_name = &lt;nil&gt;
			need_reboot = &lt;nil&gt;
			description = &lt;nil&gt;
			param_type = &lt;nil&gt;
		}

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
				Type:        schema.TypeInt,
				Description: "Page offset.",
			},

			"limit": {
				Optional:    true,
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

			"total_count": {
				Type:        schema.TypeInt,
				Description: "Number of parameter templates.",
			},

			"items": {
				Type:        schema.TypeList,
				Description: "Parameter Template Information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Description: "The ID of template.",
						},
						"template_name": {
							Type:        schema.TypeString,
							Description: "The name of template.",
						},
						"template_description": {
							Type:        schema.TypeString,
							Description: "The description of template.",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Description: "Engine version.",
						},
						"db_mode": {
							Type:        schema.TypeString,
							Description: "Database mode, optional values: NORMAL, SERVERLESS.",
						},
						"param_info_set": {
							Type:        schema.TypeList,
							Description: "Parameter template details.Note: This field may return null, indicating that no valid value can be obtained.",
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
										Description: "An optional set of value types when the parameter type is enum.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"max": {
										Type:        schema.TypeString,
										Description: "The maximum value when the parameter type is float/integer.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"min": {
										Type:        schema.TypeString,
										Description: "The minimum value when the parameter type is float/integer.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"param_name": {
										Type:        schema.TypeString,
										Description: "The name of parameter.",
									},
									"need_reboot": {
										Type:        schema.TypeInt,
										Description: "Wheter to reboot.",
									},
									"description": {
										Type:        schema.TypeString,
										Description: "The description of parameter.",
									},
									"param_type": {
										Type:        schema.TypeString,
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
		paramMap["EngineVersions"] = helper.InterfacesStringsPoint(engineVersionsSet)
	}

	if v, ok := d.GetOk("template_names"); ok {
		templateNamesSet := v.(*schema.Set).List()
		paramMap["TemplateNames"] = helper.InterfacesStringsPoint(templateNamesSet)
	}

	if v, ok := d.GetOk("template_ids"); ok {
		templateIdsSet := v.(*schema.Set).List()
		for i := range templateIdsSet {
			templateIds := templateIdsSet[i].(int)
			paramMap["TemplateIds"] = append(paramMap["TemplateIds"], helper.IntInt64(templateIds))
		}
	}

	if v, ok := d.GetOk("db_modes"); ok {
		dbModesSet := v.(*schema.Set).List()
		paramMap["DbModes"] = helper.InterfacesStringsPoint(dbModesSet)
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("products"); ok {
		productsSet := v.(*schema.Set).List()
		paramMap["Products"] = helper.InterfacesStringsPoint(productsSet)
	}

	if v, ok := d.GetOk("template_types"); ok {
		templateTypesSet := v.(*schema.Set).List()
		paramMap["TemplateTypes"] = helper.InterfacesStringsPoint(templateTypesSet)
	}

	if v, ok := d.GetOk("engine_types"); ok {
		engineTypesSet := v.(*schema.Set).List()
		paramMap["EngineTypes"] = helper.InterfacesStringsPoint(engineTypesSet)
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_direction"); ok {
		paramMap["OrderDirection"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("total_count"); v != nil {
		paramMap["TotalCount"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("items"); ok {
		itemsSet := v.([]interface{})
		tmpSet := make([]*cynosdb.ParamTemplateListInfo, 0, len(itemsSet))

		for _, item := range itemsSet {
			paramTemplateListInfo := cynosdb.ParamTemplateListInfo{}
			paramTemplateListInfoMap := item.(map[string]interface{})

			if v, ok := paramTemplateListInfoMap["id"]; ok {
				paramTemplateListInfo.Id = helper.IntInt64(v.(int))
			}
			if v, ok := paramTemplateListInfoMap["template_name"]; ok {
				paramTemplateListInfo.TemplateName = helper.String(v.(string))
			}
			if v, ok := paramTemplateListInfoMap["template_description"]; ok {
				paramTemplateListInfo.TemplateDescription = helper.String(v.(string))
			}
			if v, ok := paramTemplateListInfoMap["engine_version"]; ok {
				paramTemplateListInfo.EngineVersion = helper.String(v.(string))
			}
			if v, ok := paramTemplateListInfoMap["db_mode"]; ok {
				paramTemplateListInfo.DbMode = helper.String(v.(string))
			}
			if v, ok := paramTemplateListInfoMap["param_info_set"]; ok {
				for _, item := range v.([]interface{}) {
					paramInfoSetMap := item.(map[string]interface{})
					templateParamInfo := cynosdb.TemplateParamInfo{}
					if v, ok := paramInfoSetMap["current_value"]; ok {
						templateParamInfo.CurrentValue = helper.String(v.(string))
					}
					if v, ok := paramInfoSetMap["default"]; ok {
						templateParamInfo.Default = helper.String(v.(string))
					}
					if v, ok := paramInfoSetMap["enum_value"]; ok {
						enumValueSet := v.(*schema.Set).List()
						templateParamInfo.EnumValue = helper.InterfacesStringsPoint(enumValueSet)
					}
					if v, ok := paramInfoSetMap["max"]; ok {
						templateParamInfo.Max = helper.String(v.(string))
					}
					if v, ok := paramInfoSetMap["min"]; ok {
						templateParamInfo.Min = helper.String(v.(string))
					}
					if v, ok := paramInfoSetMap["param_name"]; ok {
						templateParamInfo.ParamName = helper.String(v.(string))
					}
					if v, ok := paramInfoSetMap["need_reboot"]; ok {
						templateParamInfo.NeedReboot = helper.IntInt64(v.(int))
					}
					if v, ok := paramInfoSetMap["description"]; ok {
						templateParamInfo.Description = helper.String(v.(string))
					}
					if v, ok := paramInfoSetMap["param_type"]; ok {
						templateParamInfo.ParamType = helper.String(v.(string))
					}
					paramTemplateListInfo.ParamInfoSet = append(paramTemplateListInfo.ParamInfoSet, &templateParamInfo)
				}
			}
			tmpSet = append(tmpSet, &paramTemplateListInfo)
		}
		paramMap["items"] = tmpSet
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

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
