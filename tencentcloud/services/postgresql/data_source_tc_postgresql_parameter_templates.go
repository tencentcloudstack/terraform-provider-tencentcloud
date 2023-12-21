package postgresql

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudPostgresqlParameterTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresqlParameterTemplatesRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter conditions. Valid values:TemplateName, TemplateId, DBMajorVersion, DBEngine.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filter name.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "One or more filter values.",
						},
					},
				},
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting metric. Valid values:CreateTime, TemplateName, DBMajorVersion.",
			},

			"order_by_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting order. Valid values:asc (ascending order),desc (descending order).",
			},

			"list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "list of parameter templates.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "parameter template ID.",
						},
						"template_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "parameter template name.",
						},
						"db_major_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the database version to which the parameter template applies.",
						},
						"db_engine": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the database engine for which the parameter template applies.",
						},
						"template_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "parameter template description.",
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

func dataSourceTencentCloudPostgresqlParameterTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_postgresql_parameter_templates.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*postgresql.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := postgresql.Filter{}
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
		paramMap["filters"] = tmpSet
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by_type"); ok {
		paramMap["OrderByType"] = helper.String(v.(string))
	}

	service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var parameterTemplateSet []*postgresql.ParameterTemplate

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribePostgresqlParameterTemplatesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		parameterTemplateSet = result
		return nil
	})
	if err != nil {
		return err
	}
	ids := make([]string, 0, len(parameterTemplateSet))
	tmpList := make([]map[string]interface{}, 0, len(parameterTemplateSet))

	if parameterTemplateSet != nil {
		for _, parameterTemplate := range parameterTemplateSet {
			parameterTemplateMap := map[string]interface{}{}

			if parameterTemplate.TemplateId != nil {
				parameterTemplateMap["template_id"] = parameterTemplate.TemplateId
			}

			if parameterTemplate.TemplateName != nil {
				parameterTemplateMap["template_name"] = parameterTemplate.TemplateName
			}

			if parameterTemplate.DBMajorVersion != nil {
				parameterTemplateMap["db_major_version"] = parameterTemplate.DBMajorVersion
			}

			if parameterTemplate.DBEngine != nil {
				parameterTemplateMap["db_engine"] = parameterTemplate.DBEngine
			}

			if parameterTemplate.TemplateDescription != nil {
				parameterTemplateMap["template_description"] = parameterTemplate.TemplateDescription
			}

			ids = append(ids, *parameterTemplate.TemplateId)
			tmpList = append(tmpList, parameterTemplateMap)
		}

		_ = d.Set("list", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
