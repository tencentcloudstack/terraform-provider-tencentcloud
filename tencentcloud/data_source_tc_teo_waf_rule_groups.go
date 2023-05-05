/*
Use this data source to query detailed information of teo wafRuleGroups

Example Usage

```hcl
data "tencentcloud_teo_waf_rule_groups" "wafRuleGroups" {
  zone_id = ""
  entity = ""
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
)

func dataSourceTencentCloudTeoWafRuleGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoWafRuleGroupsRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID.",
			},

			"entity": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subdomain or application name.",
			},

			"waf_rule_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of WAF rule groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_type_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Type id of rules in this group.",
						},
						"rule_type_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type name of rules in this group.",
						},
						"rule_type_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of rule type in this group.",
						},
						"rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Rules detail.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "WAF managed rule id.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description of the rule.",
									},
									"rule_level_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "System default level of the rule.",
									},
									"rule_tags": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Tags of the rule. Note: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudTeoWafRuleGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_teo_waf_rule_groups.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId  = getLogId(contextNil)
		ctx    = context.WithValue(context.TODO(), logIdKey, logId)
		zoneId string
		entity string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		paramMap["zone_id"] = v
	}

	if v, ok := d.GetOk("entity"); ok {
		entity = v.(string)
		paramMap["entity"] = v
	}

	teoService := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	var wafGroupDetails []*teo.WafGroupDetail
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := teoService.DescribeTeoWafRuleGroupsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		wafGroupDetails = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Teo wafGroupInfo failed, reason:%+v", logId, err)
		return err
	}

	wafGroupInfoList := []interface{}{}
	if wafGroupDetails != nil {
		for _, wafGroupDetail := range wafGroupDetails {
			wafGroupInfoMap := map[string]interface{}{}
			if wafGroupDetail.RuleTypeId != nil {
				wafGroupInfoMap["rule_type_id"] = wafGroupDetail.RuleTypeId
			}
			if wafGroupDetail.RuleTypeName != nil {
				wafGroupInfoMap["rule_type_name"] = wafGroupDetail.RuleTypeName
			}
			if wafGroupDetail.RuleTypeDesc != nil {
				wafGroupInfoMap["rule_type_desc"] = wafGroupDetail.RuleTypeDesc
			}
			if wafGroupDetail.WafGroupRules != nil {
				rulesList := []interface{}{}
				for _, rules := range wafGroupDetail.WafGroupRules {
					rulesMap := map[string]interface{}{}
					if rules.RuleId != nil {
						rulesMap["rule_id"] = rules.RuleId
					}
					if rules.Description != nil {
						rulesMap["description"] = rules.Description
					}
					if rules.RuleLevelDesc != nil {
						rulesMap["rule_level_desc"] = rules.RuleLevelDesc
					}
					if rules.RuleTags != nil {
						rulesMap["rule_tags"] = rules.RuleTags
					}

					rulesList = append(rulesList, rulesMap)
				}
				wafGroupInfoMap["rules"] = rulesList
			}

			wafGroupInfoList = append(wafGroupInfoList, wafGroupInfoMap)
		}
		_ = d.Set("waf_rule_groups", wafGroupInfoList)
	}

	d.SetId(zoneId + FILED_SP + entity)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), wafGroupInfoList); e != nil {
			return e
		}
	}
	return nil
}
