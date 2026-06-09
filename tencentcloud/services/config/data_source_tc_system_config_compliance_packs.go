package config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSystemConfigCompliancePacks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSystemConfigCompliancePacksRead,
		Schema: map[string]*schema.Schema{
			"compliance_pack_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "System compliance pack list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compliance_pack_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Compliance pack ID.",
						},
						"compliance_pack_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Compliance pack name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Compliance pack description.",
						},
						"risk_level": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Risk level. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).",
						},
						"config_rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Config rules in the compliance pack.",
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
										Description: "Rule risk level. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule creation time.",
									},
									"update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule last update time.",
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

func dataSourceTencentCloudSystemConfigCompliancePacksRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_system_config_compliance_packs.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ConfigService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	respData, reqErr := service.DescribeSystemConfigCompliancePacks(ctx)
	if reqErr != nil {
		return reqErr
	}

	packList := flattenSystemCompliancePackList(respData)
	_ = d.Set("compliance_pack_list", packList)

	d.SetId(helper.BuildToken())

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}

func flattenSystemCompliancePackList(items []*configv20220802.SystemCompliancePack) []map[string]interface{} {
	packList := make([]map[string]interface{}, 0, len(items))
	for _, pack := range items {
		packMap := map[string]interface{}{}

		if pack.CompliancePackId != nil {
			packMap["compliance_pack_id"] = pack.CompliancePackId
		}

		if pack.CompliancePackName != nil {
			packMap["compliance_pack_name"] = pack.CompliancePackName
		}

		if pack.Description != nil {
			packMap["description"] = pack.Description
		}

		if pack.RiskLevel != nil {
			packMap["risk_level"] = int(*pack.RiskLevel)
		}

		packMap["config_rules"] = flattenCompliancePackRulesForManage(pack.ConfigRules)

		packList = append(packList, packMap)
	}

	return packList
}

func flattenCompliancePackRulesForManage(rules []*configv20220802.CompliancePackRuleForManage) []map[string]interface{} {
	ruleList := make([]map[string]interface{}, 0, len(rules))
	for _, rule := range rules {
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

		if rule.CreateTime != nil {
			ruleMap["create_time"] = rule.CreateTime
		}

		if rule.UpdateTime != nil {
			ruleMap["update_time"] = rule.UpdateTime
		}

		ruleList = append(ruleList, ruleMap)
	}

	return ruleList
}
