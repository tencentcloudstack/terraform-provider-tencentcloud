/*
Provides a resource to create a wedata dq_rule

Example Usage

```hcl
resource "tencentcloud_wedata_dq_rule" "dq_rule" {
  project_id = ""
  rule_group_id =
  name = ""
  table_id = ""
  rule_template_id =
  type =
  quality_dim =
  source_object_data_type_name = ""
  source_object_value = ""
  condition_type =
  condition_expression = ""
  custom_sql = ""
  compare_rule {
		items {
			compare_type =
			operator = ""
			value_compute_type =
			value_list {
				value_type =
				value = ""
			}
		}
		cycle_step =

  }
  alarm_level =
  description = ""
  datasource_id = ""
  database_id = ""
  target_database_id = ""
  target_table_id = ""
  target_condition_expr = ""
  rel_condition_expr = ""
  field_config {
		where_config {
			field_key = ""
			field_value = ""
			field_data_type = ""
		}
		table_config {
			database_id = ""
			database_name = ""
			table_id = ""
			table_name = ""
			table_key = ""
			field_config {
				field_key = ""
				field_value = ""
				field_data_type = ""
			}
		}

  }
  target_object_value = ""
  source_engine_types =
}
```

Import

wedata dq_rule can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_dq_rule.dq_rule dq_rule_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudWedataDq_rule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataDq_ruleCreate,
		Read:   resourceTencentCloudWedataDq_ruleRead,
		Update: resourceTencentCloudWedataDq_ruleUpdate,
		Delete: resourceTencentCloudWedataDq_ruleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Project id.",
			},

			"rule_group_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Rule group id.",
			},

			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Rule name.",
			},

			"table_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Table id.",
			},

			"rule_template_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Rule template id.",
			},

			"type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Rule Type 1. System Template, 2. Custom Template, 3. Custom SQL.",
			},

			"quality_dim": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Rules belong to quality dimensions (1. accuracy, 2. uniqueness, 3. completeness, 4. consistency, 5. timeliness, 6. effectiveness).",
			},

			"source_object_data_type_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Source field type，int、string.",
			},

			"source_object_value": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Source field name.",
			},

			"condition_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Detection scope 1. Full Table 2. Conditional scan.",
			},

			"condition_expression": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Condition scans WHERE condition expressions.",
			},

			"custom_sql": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Custom sql.",
			},

			"compare_rule": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Alarm trigger condition.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"items": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Comparison condition listNote: This field may return null, indicating that a valid value cannot be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"compare_type": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Comparison type 1. Fixed value 2. Fluctuating value 3. Comparison of value range 4. Enumeration range comparison 5. Do not compareNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"operator": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Comparison operation type &amp;lt; &amp;lt;= == =&amp;gt; &amp;gt;Note: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"value_compute_type": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Quality statistics Type 1. Absolute value 2. Increase 3. Decrease 4. C contains 5. N C does not containNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"value_list": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Compare the threshold listNote: This field may return null, indicating that a valid value cannot be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value_type": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Threshold type 1. Low threshold 2. High threshold 3. Common threshold 4. Enumerated valueNote: This field may return null, indicating that a valid value cannot be obtained.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Threshold valueNote: This field may return null, indicating that a valid value cannot be obtained.",
												},
											},
										},
									},
								},
							},
						},
						"cycle_step": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Periodic Indicates the default period of a template, in secondsNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
					},
				},
			},

			"alarm_level": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Alarm trigger levels 1. Low, 2. Medium, 3. High.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Rule description.",
			},

			"datasource_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Datasource id.",
			},

			"database_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Database id.",
			},

			"target_database_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Target database id.",
			},

			"target_table_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Target table id.",
			},

			"target_condition_expr": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Target filter condition expression.",
			},

			"rel_condition_expr": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The source field and the target field are associated with a conditional on expression.",
			},

			"field_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Custom template sql expression field replacement parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"where_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Where variableNote: This field may return null, indicating that a valid value cannot be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Field keyNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"field_value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Field valueNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"field_data_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Field typeNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
								},
							},
						},
						"table_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Library table variableNote: This field may return null, indicating that a valid value cannot be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Database idNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"database_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Database nameNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"table_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Table idNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"table_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Table nameNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"table_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Table keyNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"field_config": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Field variableNote: This field may return null, indicating that a valid value cannot be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field_key": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Field keyNote: This field may return null, indicating that a valid value cannot be obtained.",
												},
												"field_value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Field valueNote: This field may return null, indicating that a valid value cannot be obtained.",
												},
												"field_data_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Field typeNote: This field may return null, indicating that a valid value cannot be obtained.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"target_object_value": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Target field name  CITY.",
			},

			"source_engine_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "List of execution engines supported by this rule.",
			},
		},
	}
}

func resourceTencentCloudWedataDq_ruleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_dq_rule.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = wedata.NewCreateRuleRequest()
		response = wedata.NewCreateRuleResponse()
		ruleId   int
	)
	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("rule_group_id"); ok {
		request.RuleGroupId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("table_id"); ok {
		request.TableId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("rule_template_id"); ok {
		request.RuleTemplateId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("type"); ok {
		request.Type = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("quality_dim"); ok {
		request.QualityDim = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("source_object_data_type_name"); ok {
		request.SourceObjectDataTypeName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_object_value"); ok {
		request.SourceObjectValue = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("condition_type"); ok {
		request.ConditionType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("condition_expression"); ok {
		request.ConditionExpression = helper.String(v.(string))
	}

	if v, ok := d.GetOk("custom_sql"); ok {
		request.CustomSql = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "compare_rule"); ok {
		compareRule := wedata.CompareRule{}
		if v, ok := dMap["items"]; ok {
			for _, item := range v.([]interface{}) {
				itemsMap := item.(map[string]interface{})
				compareRuleItem := wedata.CompareRuleItem{}
				if v, ok := itemsMap["compare_type"]; ok {
					compareRuleItem.CompareType = helper.IntUint64(v.(int))
				}
				if v, ok := itemsMap["operator"]; ok {
					compareRuleItem.Operator = helper.String(v.(string))
				}
				if v, ok := itemsMap["value_compute_type"]; ok {
					compareRuleItem.ValueComputeType = helper.IntUint64(v.(int))
				}
				if v, ok := itemsMap["value_list"]; ok {
					for _, item := range v.([]interface{}) {
						valueListMap := item.(map[string]interface{})
						thresholdValue := wedata.ThresholdValue{}
						if v, ok := valueListMap["value_type"]; ok {
							thresholdValue.ValueType = helper.IntUint64(v.(int))
						}
						if v, ok := valueListMap["value"]; ok {
							thresholdValue.Value = helper.String(v.(string))
						}
						compareRuleItem.ValueList = append(compareRuleItem.ValueList, &thresholdValue)
					}
				}
				compareRule.Items = append(compareRule.Items, &compareRuleItem)
			}
		}
		if v, ok := dMap["cycle_step"]; ok {
			compareRule.CycleStep = helper.IntUint64(v.(int))
		}
		request.CompareRule = &compareRule
	}

	if v, ok := d.GetOkExists("alarm_level"); ok {
		request.AlarmLevel = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("datasource_id"); ok {
		request.DatasourceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("database_id"); ok {
		request.DatabaseId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_database_id"); ok {
		request.TargetDatabaseId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_table_id"); ok {
		request.TargetTableId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_condition_expr"); ok {
		request.TargetConditionExpr = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rel_condition_expr"); ok {
		request.RelConditionExpr = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "field_config"); ok {
		ruleFieldConfig := wedata.RuleFieldConfig{}
		if v, ok := dMap["where_config"]; ok {
			for _, item := range v.([]interface{}) {
				whereConfigMap := item.(map[string]interface{})
				fieldConfig := wedata.FieldConfig{}
				if v, ok := whereConfigMap["field_key"]; ok {
					fieldConfig.FieldKey = helper.String(v.(string))
				}
				if v, ok := whereConfigMap["field_value"]; ok {
					fieldConfig.FieldValue = helper.String(v.(string))
				}
				if v, ok := whereConfigMap["field_data_type"]; ok {
					fieldConfig.FieldDataType = helper.String(v.(string))
				}
				ruleFieldConfig.WhereConfig = append(ruleFieldConfig.WhereConfig, &fieldConfig)
			}
		}
		if v, ok := dMap["table_config"]; ok {
			for _, item := range v.([]interface{}) {
				tableConfigMap := item.(map[string]interface{})
				tableConfig := wedata.TableConfig{}
				if v, ok := tableConfigMap["database_id"]; ok {
					tableConfig.DatabaseId = helper.String(v.(string))
				}
				if v, ok := tableConfigMap["database_name"]; ok {
					tableConfig.DatabaseName = helper.String(v.(string))
				}
				if v, ok := tableConfigMap["table_id"]; ok {
					tableConfig.TableId = helper.String(v.(string))
				}
				if v, ok := tableConfigMap["table_name"]; ok {
					tableConfig.TableName = helper.String(v.(string))
				}
				if v, ok := tableConfigMap["table_key"]; ok {
					tableConfig.TableKey = helper.String(v.(string))
				}
				if v, ok := tableConfigMap["field_config"]; ok {
					for _, item := range v.([]interface{}) {
						fieldConfigMap := item.(map[string]interface{})
						fieldConfig := wedata.FieldConfig{}
						if v, ok := fieldConfigMap["field_key"]; ok {
							fieldConfig.FieldKey = helper.String(v.(string))
						}
						if v, ok := fieldConfigMap["field_value"]; ok {
							fieldConfig.FieldValue = helper.String(v.(string))
						}
						if v, ok := fieldConfigMap["field_data_type"]; ok {
							fieldConfig.FieldDataType = helper.String(v.(string))
						}
						tableConfig.FieldConfig = append(tableConfig.FieldConfig, &fieldConfig)
					}
				}
				ruleFieldConfig.TableConfig = append(ruleFieldConfig.TableConfig, &tableConfig)
			}
		}
		request.FieldConfig = &ruleFieldConfig
	}

	if v, ok := d.GetOk("target_object_value"); ok {
		request.TargetObjectValue = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_engine_types"); ok {
		sourceEngineTypesSet := v.(*schema.Set).List()
		for i := range sourceEngineTypesSet {
			sourceEngineTypes := sourceEngineTypesSet[i].(int)
			request.SourceEngineTypes = append(request.SourceEngineTypes, helper.IntUint64(sourceEngineTypes))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().CreateRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create wedata dq_rule failed, reason:%+v", logId, err)
		return err
	}

	ruleId = *response.Response.RuleId
	d.SetId(helper.Int64ToStr(int64(ruleId)))

	return resourceTencentCloudWedataDq_ruleRead(d, meta)
}

func resourceTencentCloudWedataDq_ruleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_dq_rule.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := WedataService{client: meta.(*TencentCloudClient).apiV3Conn}

	dq_ruleId := d.Id()

	dq_rule, err := service.DescribeWedataDq_ruleById(ctx, ruleId)
	if err != nil {
		return err
	}

	if dq_rule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WedataDq_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if dq_rule.ProjectId != nil {
		_ = d.Set("project_id", dq_rule.ProjectId)
	}

	if dq_rule.RuleGroupId != nil {
		_ = d.Set("rule_group_id", dq_rule.RuleGroupId)
	}

	if dq_rule.Name != nil {
		_ = d.Set("name", dq_rule.Name)
	}

	if dq_rule.TableId != nil {
		_ = d.Set("table_id", dq_rule.TableId)
	}

	if dq_rule.RuleTemplateId != nil {
		_ = d.Set("rule_template_id", dq_rule.RuleTemplateId)
	}

	if dq_rule.Type != nil {
		_ = d.Set("type", dq_rule.Type)
	}

	if dq_rule.QualityDim != nil {
		_ = d.Set("quality_dim", dq_rule.QualityDim)
	}

	if dq_rule.SourceObjectDataTypeName != nil {
		_ = d.Set("source_object_data_type_name", dq_rule.SourceObjectDataTypeName)
	}

	if dq_rule.SourceObjectValue != nil {
		_ = d.Set("source_object_value", dq_rule.SourceObjectValue)
	}

	if dq_rule.ConditionType != nil {
		_ = d.Set("condition_type", dq_rule.ConditionType)
	}

	if dq_rule.ConditionExpression != nil {
		_ = d.Set("condition_expression", dq_rule.ConditionExpression)
	}

	if dq_rule.CustomSql != nil {
		_ = d.Set("custom_sql", dq_rule.CustomSql)
	}

	if dq_rule.CompareRule != nil {
		compareRuleMap := map[string]interface{}{}

		if dq_rule.CompareRule.Items != nil {
			itemsList := []interface{}{}
			for _, items := range dq_rule.CompareRule.Items {
				itemsMap := map[string]interface{}{}

				if items.CompareType != nil {
					itemsMap["compare_type"] = items.CompareType
				}

				if items.Operator != nil {
					itemsMap["operator"] = items.Operator
				}

				if items.ValueComputeType != nil {
					itemsMap["value_compute_type"] = items.ValueComputeType
				}

				if items.ValueList != nil {
					valueListList := []interface{}{}
					for _, valueList := range items.ValueList {
						valueListMap := map[string]interface{}{}

						if valueList.ValueType != nil {
							valueListMap["value_type"] = valueList.ValueType
						}

						if valueList.Value != nil {
							valueListMap["value"] = valueList.Value
						}

						valueListList = append(valueListList, valueListMap)
					}

					itemsMap["value_list"] = []interface{}{valueListList}
				}

				itemsList = append(itemsList, itemsMap)
			}

			compareRuleMap["items"] = []interface{}{itemsList}
		}

		if dq_rule.CompareRule.CycleStep != nil {
			compareRuleMap["cycle_step"] = dq_rule.CompareRule.CycleStep
		}

		_ = d.Set("compare_rule", []interface{}{compareRuleMap})
	}

	if dq_rule.AlarmLevel != nil {
		_ = d.Set("alarm_level", dq_rule.AlarmLevel)
	}

	if dq_rule.Description != nil {
		_ = d.Set("description", dq_rule.Description)
	}

	if dq_rule.DatasourceId != nil {
		_ = d.Set("datasource_id", dq_rule.DatasourceId)
	}

	if dq_rule.DatabaseId != nil {
		_ = d.Set("database_id", dq_rule.DatabaseId)
	}

	if dq_rule.TargetDatabaseId != nil {
		_ = d.Set("target_database_id", dq_rule.TargetDatabaseId)
	}

	if dq_rule.TargetTableId != nil {
		_ = d.Set("target_table_id", dq_rule.TargetTableId)
	}

	if dq_rule.TargetConditionExpr != nil {
		_ = d.Set("target_condition_expr", dq_rule.TargetConditionExpr)
	}

	if dq_rule.RelConditionExpr != nil {
		_ = d.Set("rel_condition_expr", dq_rule.RelConditionExpr)
	}

	if dq_rule.FieldConfig != nil {
		fieldConfigMap := map[string]interface{}{}

		if dq_rule.FieldConfig.WhereConfig != nil {
			whereConfigList := []interface{}{}
			for _, whereConfig := range dq_rule.FieldConfig.WhereConfig {
				whereConfigMap := map[string]interface{}{}

				if whereConfig.FieldKey != nil {
					whereConfigMap["field_key"] = whereConfig.FieldKey
				}

				if whereConfig.FieldValue != nil {
					whereConfigMap["field_value"] = whereConfig.FieldValue
				}

				if whereConfig.FieldDataType != nil {
					whereConfigMap["field_data_type"] = whereConfig.FieldDataType
				}

				whereConfigList = append(whereConfigList, whereConfigMap)
			}

			fieldConfigMap["where_config"] = []interface{}{whereConfigList}
		}

		if dq_rule.FieldConfig.TableConfig != nil {
			tableConfigList := []interface{}{}
			for _, tableConfig := range dq_rule.FieldConfig.TableConfig {
				tableConfigMap := map[string]interface{}{}

				if tableConfig.DatabaseId != nil {
					tableConfigMap["database_id"] = tableConfig.DatabaseId
				}

				if tableConfig.DatabaseName != nil {
					tableConfigMap["database_name"] = tableConfig.DatabaseName
				}

				if tableConfig.TableId != nil {
					tableConfigMap["table_id"] = tableConfig.TableId
				}

				if tableConfig.TableName != nil {
					tableConfigMap["table_name"] = tableConfig.TableName
				}

				if tableConfig.TableKey != nil {
					tableConfigMap["table_key"] = tableConfig.TableKey
				}

				if tableConfig.FieldConfig != nil {
					fieldConfigList := []interface{}{}
					for _, fieldConfig := range tableConfig.FieldConfig {
						fieldConfigMap := map[string]interface{}{}

						if fieldConfig.FieldKey != nil {
							fieldConfigMap["field_key"] = fieldConfig.FieldKey
						}

						if fieldConfig.FieldValue != nil {
							fieldConfigMap["field_value"] = fieldConfig.FieldValue
						}

						if fieldConfig.FieldDataType != nil {
							fieldConfigMap["field_data_type"] = fieldConfig.FieldDataType
						}

						fieldConfigList = append(fieldConfigList, fieldConfigMap)
					}

					tableConfigMap["field_config"] = []interface{}{fieldConfigList}
				}

				tableConfigList = append(tableConfigList, tableConfigMap)
			}

			fieldConfigMap["table_config"] = []interface{}{tableConfigList}
		}

		_ = d.Set("field_config", []interface{}{fieldConfigMap})
	}

	if dq_rule.TargetObjectValue != nil {
		_ = d.Set("target_object_value", dq_rule.TargetObjectValue)
	}

	if dq_rule.SourceEngineTypes != nil {
		_ = d.Set("source_engine_types", dq_rule.SourceEngineTypes)
	}

	return nil
}

func resourceTencentCloudWedataDq_ruleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_dq_rule.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := wedata.NewModifyRuleRequest()

	dq_ruleId := d.Id()

	request.RuleId = &ruleId

	immutableArgs := []string{"project_id", "rule_group_id", "name", "table_id", "rule_template_id", "type", "quality_dim", "source_object_data_type_name", "source_object_value", "condition_type", "condition_expression", "custom_sql", "compare_rule", "alarm_level", "description", "datasource_id", "database_id", "target_database_id", "target_table_id", "target_condition_expr", "rel_condition_expr", "field_config", "target_object_value", "source_engine_types"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("project_id") {
		if v, ok := d.GetOk("project_id"); ok {
			request.ProjectId = helper.String(v.(string))
		}
	}

	if d.HasChange("rule_group_id") {
		if v, ok := d.GetOkExists("rule_group_id"); ok {
			request.RuleGroupId = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("table_id") {
		if v, ok := d.GetOk("table_id"); ok {
			request.TableId = helper.String(v.(string))
		}
	}

	if d.HasChange("rule_template_id") {
		if v, ok := d.GetOkExists("rule_template_id"); ok {
			request.RuleTemplateId = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("type") {
		if v, ok := d.GetOkExists("type"); ok {
			request.Type = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("quality_dim") {
		if v, ok := d.GetOkExists("quality_dim"); ok {
			request.QualityDim = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("source_object_data_type_name") {
		if v, ok := d.GetOk("source_object_data_type_name"); ok {
			request.SourceObjectDataTypeName = helper.String(v.(string))
		}
	}

	if d.HasChange("source_object_value") {
		if v, ok := d.GetOk("source_object_value"); ok {
			request.SourceObjectValue = helper.String(v.(string))
		}
	}

	if d.HasChange("condition_type") {
		if v, ok := d.GetOkExists("condition_type"); ok {
			request.ConditionType = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("condition_expression") {
		if v, ok := d.GetOk("condition_expression"); ok {
			request.ConditionExpression = helper.String(v.(string))
		}
	}

	if d.HasChange("custom_sql") {
		if v, ok := d.GetOk("custom_sql"); ok {
			request.CustomSql = helper.String(v.(string))
		}
	}

	if d.HasChange("compare_rule") {
		if dMap, ok := helper.InterfacesHeadMap(d, "compare_rule"); ok {
			compareRule := wedata.CompareRule{}
			if v, ok := dMap["items"]; ok {
				for _, item := range v.([]interface{}) {
					itemsMap := item.(map[string]interface{})
					compareRuleItem := wedata.CompareRuleItem{}
					if v, ok := itemsMap["compare_type"]; ok {
						compareRuleItem.CompareType = helper.IntUint64(v.(int))
					}
					if v, ok := itemsMap["operator"]; ok {
						compareRuleItem.Operator = helper.String(v.(string))
					}
					if v, ok := itemsMap["value_compute_type"]; ok {
						compareRuleItem.ValueComputeType = helper.IntUint64(v.(int))
					}
					if v, ok := itemsMap["value_list"]; ok {
						for _, item := range v.([]interface{}) {
							valueListMap := item.(map[string]interface{})
							thresholdValue := wedata.ThresholdValue{}
							if v, ok := valueListMap["value_type"]; ok {
								thresholdValue.ValueType = helper.IntUint64(v.(int))
							}
							if v, ok := valueListMap["value"]; ok {
								thresholdValue.Value = helper.String(v.(string))
							}
							compareRuleItem.ValueList = append(compareRuleItem.ValueList, &thresholdValue)
						}
					}
					compareRule.Items = append(compareRule.Items, &compareRuleItem)
				}
			}
			if v, ok := dMap["cycle_step"]; ok {
				compareRule.CycleStep = helper.IntUint64(v.(int))
			}
			request.CompareRule = &compareRule
		}
	}

	if d.HasChange("alarm_level") {
		if v, ok := d.GetOkExists("alarm_level"); ok {
			request.AlarmLevel = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("target_database_id") {
		if v, ok := d.GetOk("target_database_id"); ok {
			request.TargetDatabaseId = helper.String(v.(string))
		}
	}

	if d.HasChange("target_table_id") {
		if v, ok := d.GetOk("target_table_id"); ok {
			request.TargetTableId = helper.String(v.(string))
		}
	}

	if d.HasChange("target_condition_expr") {
		if v, ok := d.GetOk("target_condition_expr"); ok {
			request.TargetConditionExpr = helper.String(v.(string))
		}
	}

	if d.HasChange("rel_condition_expr") {
		if v, ok := d.GetOk("rel_condition_expr"); ok {
			request.RelConditionExpr = helper.String(v.(string))
		}
	}

	if d.HasChange("field_config") {
		if dMap, ok := helper.InterfacesHeadMap(d, "field_config"); ok {
			ruleFieldConfig := wedata.RuleFieldConfig{}
			if v, ok := dMap["where_config"]; ok {
				for _, item := range v.([]interface{}) {
					whereConfigMap := item.(map[string]interface{})
					fieldConfig := wedata.FieldConfig{}
					if v, ok := whereConfigMap["field_key"]; ok {
						fieldConfig.FieldKey = helper.String(v.(string))
					}
					if v, ok := whereConfigMap["field_value"]; ok {
						fieldConfig.FieldValue = helper.String(v.(string))
					}
					if v, ok := whereConfigMap["field_data_type"]; ok {
						fieldConfig.FieldDataType = helper.String(v.(string))
					}
					ruleFieldConfig.WhereConfig = append(ruleFieldConfig.WhereConfig, &fieldConfig)
				}
			}
			if v, ok := dMap["table_config"]; ok {
				for _, item := range v.([]interface{}) {
					tableConfigMap := item.(map[string]interface{})
					tableConfig := wedata.TableConfig{}
					if v, ok := tableConfigMap["database_id"]; ok {
						tableConfig.DatabaseId = helper.String(v.(string))
					}
					if v, ok := tableConfigMap["database_name"]; ok {
						tableConfig.DatabaseName = helper.String(v.(string))
					}
					if v, ok := tableConfigMap["table_id"]; ok {
						tableConfig.TableId = helper.String(v.(string))
					}
					if v, ok := tableConfigMap["table_name"]; ok {
						tableConfig.TableName = helper.String(v.(string))
					}
					if v, ok := tableConfigMap["table_key"]; ok {
						tableConfig.TableKey = helper.String(v.(string))
					}
					if v, ok := tableConfigMap["field_config"]; ok {
						for _, item := range v.([]interface{}) {
							fieldConfigMap := item.(map[string]interface{})
							fieldConfig := wedata.FieldConfig{}
							if v, ok := fieldConfigMap["field_key"]; ok {
								fieldConfig.FieldKey = helper.String(v.(string))
							}
							if v, ok := fieldConfigMap["field_value"]; ok {
								fieldConfig.FieldValue = helper.String(v.(string))
							}
							if v, ok := fieldConfigMap["field_data_type"]; ok {
								fieldConfig.FieldDataType = helper.String(v.(string))
							}
							tableConfig.FieldConfig = append(tableConfig.FieldConfig, &fieldConfig)
						}
					}
					ruleFieldConfig.TableConfig = append(ruleFieldConfig.TableConfig, &tableConfig)
				}
			}
			request.FieldConfig = &ruleFieldConfig
		}
	}

	if d.HasChange("target_object_value") {
		if v, ok := d.GetOk("target_object_value"); ok {
			request.TargetObjectValue = helper.String(v.(string))
		}
	}

	if d.HasChange("source_engine_types") {
		if v, ok := d.GetOk("source_engine_types"); ok {
			sourceEngineTypesSet := v.(*schema.Set).List()
			for i := range sourceEngineTypesSet {
				sourceEngineTypes := sourceEngineTypesSet[i].(int)
				request.SourceEngineTypes = append(request.SourceEngineTypes, helper.IntUint64(sourceEngineTypes))
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().ModifyRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update wedata dq_rule failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWedataDq_ruleRead(d, meta)
}

func resourceTencentCloudWedataDq_ruleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_dq_rule.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := WedataService{client: meta.(*TencentCloudClient).apiV3Conn}
	dq_ruleId := d.Id()

	if err := service.DeleteWedataDq_ruleById(ctx, ruleId); err != nil {
		return err
	}

	return nil
}
