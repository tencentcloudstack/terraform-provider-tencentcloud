/*
Use this data source to query detailed information of wedata rule templates

Example Usage

```hcl
data "tencentcloud_wedata_rule_templates" "rule_templates" {
  type                = 2
  source_object_type  = 2
  project_id          = "1840731346428280832"
  source_engine_types = [2, 4, 16]
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudWedataRuleTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataRuleTemplatesRead,
		Schema: map[string]*schema.Schema{
			"type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Template type. `1` means System template, `2` means Custom template.",
			},

			"source_object_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Source data object type. `1`: Constant, `2`: Offline table level, `3`: Offline field level.",
			},

			"project_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Project ID.",
			},

			"source_engine_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Applicable type of source data.",
			},

			"data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "rule template list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_template_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of rule template.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of rule template.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of rule template.",
						},
						"type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Template type. `1` means System template, `2` means Custom template.",
						},
						"source_object_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Source object type. `1`: Constant, `2`: Offline table level, `3`: Offline field level.",
						},
						"source_object_data_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Source data object type. `1`: value, `2`: string.",
						},
						"source_content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Content of rule template.",
						},
						"source_engine_types": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "Applicable type of source data.",
						},
						"quality_dim": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Quality inspection dimensions. `1`: Accuracy, `2`: Uniqueness, `3`: Completeness, `4`: Consistency, `5`: Timeliness, `6`: Effectiveness.",
						},
						"compare_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The type of comparison method supported by the rule (1: fixed value comparison, greater than, less than, greater than or equal to, etc. 2: fluctuating value comparison, absolute value, rise, fall).",
						},
						"citation_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Citations.",
						},
						"user_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "user id.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "user name.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "update time, like: yyyy-MM-dd HH:mm:ss.",
						},
						"where_flag": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If add where.",
						},
						"multi_source_flag": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to associate other library tables.",
						},
						"sql_expression": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Sql Expression.",
						},
						"sub_quality_dim": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Sub Quality inspection dimensions. `1`: Accuracy, `2`: Uniqueness, `3`: Completeness, `4`: Consistency, `5`: Timeliness, `6`: Effectiveness.",
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

func dataSourceTencentCloudWedataRuleTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_wedata_rule_templates.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOkExists("type"); ok {
		paramMap["Type"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("source_object_type"); ok {
		paramMap["SourceObjectType"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("project_id"); ok {
		paramMap["ProjectId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_engine_types"); ok {
		sourceEngineTypesSet := v.(*schema.Set).List()
		var sourceEngineTypes []*uint64
		for i := range sourceEngineTypesSet {
			sourceEngineType := sourceEngineTypesSet[i].(int)
			sourceEngineTypes = append(sourceEngineTypes, helper.IntUint64(sourceEngineType))
		}
		paramMap["SourceEngineTypes"] = sourceEngineTypes
	}

	service := WedataService{client: meta.(*TencentCloudClient).apiV3Conn}

	var data []*wedata.RuleTemplate

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataRuleTemplatesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		data = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(data))
	tmpList := make([]map[string]interface{}, 0, len(data))

	if data != nil {
		for _, ruleTemplate := range data {
			ruleTemplateMap := map[string]interface{}{}

			if ruleTemplate.RuleTemplateId != nil {
				ruleTemplateMap["rule_template_id"] = ruleTemplate.RuleTemplateId
			}

			if ruleTemplate.Name != nil {
				ruleTemplateMap["name"] = ruleTemplate.Name
			}

			if ruleTemplate.Description != nil {
				ruleTemplateMap["description"] = ruleTemplate.Description
			}

			if ruleTemplate.Type != nil {
				ruleTemplateMap["type"] = ruleTemplate.Type
			}

			if ruleTemplate.SourceObjectType != nil {
				ruleTemplateMap["source_object_type"] = ruleTemplate.SourceObjectType
			}

			if ruleTemplate.SourceObjectDataType != nil {
				ruleTemplateMap["source_object_data_type"] = ruleTemplate.SourceObjectDataType
			}

			if ruleTemplate.SourceContent != nil {
				ruleTemplateMap["source_content"] = ruleTemplate.SourceContent
			}

			if ruleTemplate.SourceEngineTypes != nil {
				ruleTemplateMap["source_engine_types"] = ruleTemplate.SourceEngineTypes
			}

			if ruleTemplate.QualityDim != nil {
				ruleTemplateMap["quality_dim"] = ruleTemplate.QualityDim
			}

			if ruleTemplate.CompareType != nil {
				ruleTemplateMap["compare_type"] = ruleTemplate.CompareType
			}

			if ruleTemplate.CitationCount != nil {
				ruleTemplateMap["citation_count"] = ruleTemplate.CitationCount
			}

			if ruleTemplate.UserId != nil {
				ruleTemplateMap["user_id"] = ruleTemplate.UserId
			}

			if ruleTemplate.UserName != nil {
				ruleTemplateMap["user_name"] = ruleTemplate.UserName
			}

			if ruleTemplate.UpdateTime != nil {
				ruleTemplateMap["update_time"] = ruleTemplate.UpdateTime
			}

			if ruleTemplate.WhereFlag != nil {
				ruleTemplateMap["where_flag"] = ruleTemplate.WhereFlag
			}

			if ruleTemplate.MultiSourceFlag != nil {
				ruleTemplateMap["multi_source_flag"] = ruleTemplate.MultiSourceFlag
			}

			if ruleTemplate.SqlExpression != nil {
				ruleTemplateMap["sql_expression"] = ruleTemplate.SqlExpression
			}

			if ruleTemplate.SubQualityDim != nil {
				ruleTemplateMap["sub_quality_dim"] = ruleTemplate.SubQualityDim
			}

			ids = append(ids, helper.UInt64ToStr(*ruleTemplate.RuleTemplateId))
			tmpList = append(tmpList, ruleTemplateMap)
		}

		_ = d.Set("data", tmpList)
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
