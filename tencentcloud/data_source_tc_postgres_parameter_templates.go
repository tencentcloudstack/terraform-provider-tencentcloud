/*
Use this data source to query detailed information of postgres parameter_templates

Example Usage

```hcl
data "tencentcloud_postgres_parameter_templates" "parameter_templates" {
  filters {
		name = "DBEngine"
		values =

  }
  limit = 20
  offset = 0
  order_by = "CreateTime"
  order_by_type = "desc"
  template_id = "0b3eaa95-dfba-5253-8cdd-1258ae34d596"
  template_name = "test template"
  d_b_major_version = "14"
  d_b_engine = "PostgreSQL"
  template_description = "for test"
}
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudPostgresParameterTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresParameterTemplatesRead,
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

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The maximum number of results returned per page. Value range:0-100. Default:20.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Data offset.",
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

			"template_id": {
				Type:        schema.TypeString,
				Description: "Template ID.",
			},

			"template_name": {
				Type:        schema.TypeString,
				Description: "Template name.",
			},

			"d_b_major_version": {
				Type:        schema.TypeString,
				Description: "DBInstance major version.",
			},

			"d_b_engine": {
				Type:        schema.TypeString,
				Description: "DBInstance engine.",
			},

			"template_description": {
				Type:        schema.TypeString,
				Description: "Template description content.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudPostgresParameterTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_postgres_parameter_templates.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*postgres.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := postgres.Filter{}
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

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by_type"); ok {
		paramMap["OrderByType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("template_id"); ok {
		paramMap["TemplateId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("template_name"); ok {
		paramMap["TemplateName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("d_b_major_version"); ok {
		paramMap["DBMajorVersion"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("d_b_engine"); ok {
		paramMap["DBEngine"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("template_description"); ok {
		paramMap["TemplateDescription"] = helper.String(v.(string))
	}

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}

	var parameterTemplateSet []*postgres.ParameterTemplate

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribePostgresParameterTemplatesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		parameterTemplateSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(parameterTemplateSet))
	tmpList := make([]map[string]interface{}, 0, len(parameterTemplateSet))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
