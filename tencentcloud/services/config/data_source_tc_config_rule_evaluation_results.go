package config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func DataSourceTencentCloudConfigRuleEvaluationResults() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudConfigRuleEvaluationResultsRead,
		Schema: map[string]*schema.Schema{
			"config_rule_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Config rule ID.",
			},

			"resource_type": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Resource type list for filtering (e.g. QCS::CVM::Instance).",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"compliance_type": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Compliance type list for filtering. Valid values: COMPLIANT, NON_COMPLIANT.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"result_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Evaluation result list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource ID.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource type.",
						},
						"resource_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource region.",
						},
						"resource_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource name.",
						},
						"config_rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Config rule ID.",
						},
						"config_rule_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Config rule name.",
						},
						"compliance_pack_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Compliance pack ID.",
						},
						"risk_level": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Risk level. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).",
						},
						"compliance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Compliance type. Valid values: COMPLIANT, NON_COMPLIANT.",
						},
						"invoking_event_message_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule invocation type.",
						},
						"config_rule_invoked_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Evaluation invocation time.",
						},
						"result_recorded_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Evaluation result recorded time.",
						},
						"annotation": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Evaluation annotation detail.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"configuration": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Actual resource configuration (non-compliant configuration).",
									},
									"desired_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Expected resource configuration (compliant configuration).",
									},
									"operator": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Comparison operator between actual and expected configuration.",
									},
									"property": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "JSON path of the current configuration in the resource attribute structure.",
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

func dataSourceTencentCloudConfigRuleEvaluationResultsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_config_rule_evaluation_results.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ConfigService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := buildRuleEvaluationResultsParamMap(d)

	respData, reqErr := service.DescribeConfigRuleEvaluationResultsByFilter(ctx, paramMap)
	if reqErr != nil {
		return reqErr
	}

	resultList := flattenEvaluationResultList(respData)
	_ = d.Set("result_list", resultList)

	d.SetId(d.Get("config_rule_id").(string))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}

func buildRuleEvaluationResultsParamMap(d *schema.ResourceData) map[string]interface{} {
	paramMap := make(map[string]interface{})

	if v, ok := d.GetOk("config_rule_id"); ok {
		paramMap["ConfigRuleId"] = v.(string)
	}

	if v, ok := d.GetOk("resource_type"); ok {
		rawList := v.([]interface{})
		types := make([]*string, 0, len(rawList))
		for _, item := range rawList {
			val := item.(string)
			types = append(types, &val)
		}

		paramMap["ResourceType"] = types
	}

	if v, ok := d.GetOk("compliance_type"); ok {
		rawList := v.([]interface{})
		complianceTypes := make([]*string, 0, len(rawList))
		for _, item := range rawList {
			val := item.(string)
			complianceTypes = append(complianceTypes, &val)
		}

		paramMap["ComplianceType"] = complianceTypes
	}

	return paramMap
}

func flattenEvaluationResultList(items []*configv20220802.EvaluationResult) []map[string]interface{} {
	resultList := make([]map[string]interface{}, 0, len(items))
	for _, result := range items {
		resultMap := map[string]interface{}{}

		if result.ResourceId != nil {
			resultMap["resource_id"] = result.ResourceId
		}

		if result.ResourceType != nil {
			resultMap["resource_type"] = result.ResourceType
		}

		if result.ResourceRegion != nil {
			resultMap["resource_region"] = result.ResourceRegion
		}

		if result.ResourceName != nil {
			resultMap["resource_name"] = result.ResourceName
		}

		if result.ConfigRuleId != nil {
			resultMap["config_rule_id"] = result.ConfigRuleId
		}

		if result.ConfigRuleName != nil {
			resultMap["config_rule_name"] = result.ConfigRuleName
		}

		if result.CompliancePackId != nil {
			resultMap["compliance_pack_id"] = result.CompliancePackId
		}

		if result.RiskLevel != nil {
			resultMap["risk_level"] = int(*result.RiskLevel)
		}

		if result.ComplianceType != nil {
			resultMap["compliance_type"] = result.ComplianceType
		}

		if result.InvokingEventMessageType != nil {
			resultMap["invoking_event_message_type"] = result.InvokingEventMessageType
		}

		if result.ConfigRuleInvokedTime != nil {
			resultMap["config_rule_invoked_time"] = result.ConfigRuleInvokedTime
		}

		if result.ResultRecordedTime != nil {
			resultMap["result_recorded_time"] = result.ResultRecordedTime
		}

		if result.Annotation != nil {
			annotationMap := map[string]interface{}{}
			if result.Annotation.Configuration != nil {
				annotationMap["configuration"] = result.Annotation.Configuration
			}

			if result.Annotation.DesiredValue != nil {
				annotationMap["desired_value"] = result.Annotation.DesiredValue
			}

			if result.Annotation.Operator != nil {
				annotationMap["operator"] = result.Annotation.Operator
			}

			if result.Annotation.Property != nil {
				annotationMap["property"] = result.Annotation.Property
			}

			resultMap["annotation"] = []map[string]interface{}{annotationMap}
		}

		resultList = append(resultList, resultMap)
	}

	return resultList
}
