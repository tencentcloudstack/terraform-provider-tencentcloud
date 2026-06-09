package config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudConfigSystemRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudConfigSystemRulesRead,
		Schema: map[string]*schema.Schema{
			"keyword": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search keyword. Supports identifier/name/label/description search.",
			},

			"risk_level": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Risk level for filtering. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).",
			},

			"rule_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "System preset rule list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule unique identifier.",
						},
						"rule_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule description.",
						},
						"risk_level": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Risk level. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).",
						},
						"service_function": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Corresponding service function.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last update time.",
						},
						"trigger_type": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Trigger type list.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"resource_type": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Supported resource type list.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"label": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Rule label list.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"reference_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of times this rule is referenced.",
						},
						"identifier_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule type.",
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

func dataSourceTencentCloudConfigSystemRulesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_config_system_rules.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ConfigService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := buildSystemRulesParamMap(d)

	respData, reqErr := service.DescribeConfigSystemRulesByFilter(ctx, paramMap)
	if reqErr != nil {
		return reqErr
	}

	ruleList := flattenSystemConfigRuleList(respData)
	_ = d.Set("rule_list", ruleList)

	d.SetId(helper.BuildToken())

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}

func buildSystemRulesParamMap(d *schema.ResourceData) map[string]interface{} {
	paramMap := make(map[string]interface{})

	if v, ok := d.GetOk("keyword"); ok {
		paramMap["Keyword"] = v.(string)
	}

	if v, ok := d.GetOk("risk_level"); ok {
		paramMap["RiskLevel"] = v.(int)
	}

	return paramMap
}

func flattenSystemConfigRuleList(items []*configv20220802.SystemConfigRule) []map[string]interface{} {
	ruleList := make([]map[string]interface{}, 0, len(items))
	for _, rule := range items {
		ruleMap := map[string]interface{}{}

		if rule.Identifier != nil {
			ruleMap["identifier"] = rule.Identifier
		}

		if rule.RuleName != nil {
			ruleMap["rule_name"] = rule.RuleName
		}

		if rule.Description != nil {
			ruleMap["description"] = rule.Description
		}

		if rule.RiskLevel != nil {
			ruleMap["risk_level"] = int(*rule.RiskLevel)
		}

		if rule.ServiceFunction != nil {
			ruleMap["service_function"] = rule.ServiceFunction
		}

		if rule.CreateTime != nil {
			ruleMap["create_time"] = rule.CreateTime
		}

		if rule.UpdateTime != nil {
			ruleMap["update_time"] = rule.UpdateTime
		}

		if rule.TriggerType != nil {
			triggerTypes := make([]string, 0, len(rule.TriggerType))
			for _, t := range rule.TriggerType {
				if t != nil {
					triggerTypes = append(triggerTypes, *t)
				}
			}

			ruleMap["trigger_type"] = triggerTypes
		}

		if rule.ResourceType != nil {
			resourceTypes := make([]string, 0, len(rule.ResourceType))
			for _, rt := range rule.ResourceType {
				if rt != nil {
					resourceTypes = append(resourceTypes, *rt)
				}
			}

			ruleMap["resource_type"] = resourceTypes
		}

		if rule.Label != nil {
			labels := make([]string, 0, len(rule.Label))
			for _, l := range rule.Label {
				if l != nil {
					labels = append(labels, *l)
				}
			}

			ruleMap["label"] = labels
		}

		if rule.ReferenceCount != nil {
			ruleMap["reference_count"] = int(*rule.ReferenceCount)
		}

		if rule.IdentifierType != nil {
			ruleMap["identifier_type"] = rule.IdentifierType
		}

		ruleList = append(ruleList, ruleMap)
	}

	return ruleList
}
