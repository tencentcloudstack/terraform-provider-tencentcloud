package waf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wafv20180125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWafOwaspRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWafOwaspRulesRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain to be queried.",
			},

			"by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the field used to sort. valid values: RuleId, ModifyTime.",
			},

			"order": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sorting method. supports asc, desc.",
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies the criteria, support RuleId, TypeId, Desc, CveID, Status, and VulLevel.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Field name, used for filtering\nFilter the sub-order number (value) by DealName.",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Values after filtering.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"exact_match": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Exact search or not.",
						},
					},
				},
			},

			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Rule ID.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule description.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Rule switch. valid values: 0 (disabled), 1 (enabled), 2 (observation only).",
						},
						"level": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Protection level of the rule. valid values: 100 (loose), 200 (normal), 300 (strict), 400 (ultra-strict).",
						},
						"vul_level": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Threat level. valid values: 0 (unknown), 100 (low risk), 200 (medium risk), 300 (high risk), 400 (critical).",
						},
						"cve_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CVE ID.",
						},
						"type_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the rule type ID.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time.",
						},
						"locked": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether the user is locked.",
						},
						"reason": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Reason for modification\n\n0: none (compatibility records are empty).\n1: avoid false positives due to business characteristics.\n2: reporting of rule-based false positives.\n3: gray release of core business rules.\n4: others.",
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

func dataSourceTencentCloudWafOwaspRulesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_waf_owasp_rules.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		domain  string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("domain"); ok {
		paramMap["Domain"] = helper.String(v.(string))
		domain = v.(string)
	}

	if v, ok := d.GetOk("by"); ok {
		paramMap["By"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order"); ok {
		paramMap["Order"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*wafv20180125.FiltersItemNew, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filtersItemNew := wafv20180125.FiltersItemNew{}
			if v, ok := filtersMap["name"].(string); ok && v != "" {
				filtersItemNew.Name = helper.String(v)
			}

			if v, ok := filtersMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				for i := range valuesSet {
					values := valuesSet[i].(string)
					filtersItemNew.Values = append(filtersItemNew.Values, helper.String(values))
				}
			}

			if v, ok := filtersMap["exact_match"].(bool); ok {
				filtersItemNew.ExactMatch = helper.Bool(v)
			}

			tmpSet = append(tmpSet, &filtersItemNew)
		}

		paramMap["Filters"] = tmpSet
	}

	var respData []*wafv20180125.OwaspRule
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWafOwaspRulesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	listList := make([]map[string]interface{}, 0, len(respData))
	for _, list := range respData {
		listMap := map[string]interface{}{}
		if list.RuleId != nil {
			listMap["rule_id"] = list.RuleId
		}

		if list.Description != nil {
			listMap["description"] = list.Description
		}

		if list.Status != nil {
			listMap["status"] = list.Status
		}

		if list.Level != nil {
			listMap["level"] = list.Level
		}

		if list.VulLevel != nil {
			listMap["vul_level"] = list.VulLevel
		}

		if list.CveID != nil {
			listMap["cve_id"] = list.CveID
		}

		if list.TypeId != nil {
			listMap["type_id"] = list.TypeId
		}

		if list.CreateTime != nil {
			listMap["create_time"] = list.CreateTime
		}

		if list.ModifyTime != nil {
			listMap["modify_time"] = list.ModifyTime
		}

		if list.Locked != nil {
			listMap["locked"] = list.Locked
		}

		if list.Reason != nil {
			listMap["reason"] = list.Reason
		}

		listList = append(listList, listMap)
	}

	_ = d.Set("list", listList)

	d.SetId(domain)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
