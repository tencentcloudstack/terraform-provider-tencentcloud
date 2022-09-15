/*
Use this data source to query detailed information of teo wafManagedRules

Example Usage

```hcl
data "tencentcloud_teo_waf_managed_rules" "wafManagedRules" {
  zone_id = ""
  entity = ""
    }
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
)

func dataSourceTencentCloudTeoWafManagedRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoWafManagedRulesRead,
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

			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total managed rule number.",
			},

			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Managed rules list.",
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
							Description: "Description of the rule.",
						},
						"rule_type_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the rule.",
						},
						"rule_level_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protection level of the rule.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last modification date.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the rule. Valid values: `allow`, `block`.",
						},
						"rule_tags": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Tags of the rule. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"rule_type_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the rule type. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"rule_type_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Type ID of the rule. Note: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudTeoWafManagedRulesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_teo_waf_managed_rules.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone_id"); ok {
		paramMap["zone_id"] = v
	}

	if v, ok := d.GetOk("entity"); ok {
		paramMap["entity"] = v
	}

	teoService := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	var rules []*teo.ManagedRule
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := teoService.DescribeTeoWafManagedRulesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		rules = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Teo rules failed, reason:%+v", logId, err)
		return err
	}

	if rules != nil {
		_ = d.Set("total", len(rules))
	}

	ruleList := []interface{}{}
	if rules != nil {
		for _, rule := range rules {
			ruleMap := map[string]interface{}{}
			if rule.RuleId != nil {
				ruleMap["rule_id"] = rule.RuleId
			}
			if rule.Description != nil {
				ruleMap["description"] = rule.Description
			}
			if rule.RuleTypeName != nil {
				ruleMap["rule_type_name"] = rule.RuleTypeName
			}
			if rule.RuleLevelDesc != nil {
				ruleMap["rule_level_desc"] = rule.RuleLevelDesc
			}
			if rule.UpdateTime != nil {
				ruleMap["update_time"] = rule.UpdateTime
			}
			if rule.Status != nil {
				ruleMap["status"] = rule.Status
			}
			if rule.RuleTags != nil {
				ruleMap["rule_tags"] = rule.RuleTags
			}
			if rule.RuleTypeDesc != nil {
				ruleMap["rule_type_desc"] = rule.RuleTypeDesc
			}
			if rule.RuleTypeId != nil {
				ruleMap["rule_type_id"] = rule.RuleTypeId
			}

			ruleList = append(ruleList, ruleMap)
		}
		_ = d.Set("rules", ruleList)
	}

	d.SetId("waf_managed_rules")

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), ruleList); e != nil {
			return e
		}
	}
	return nil
}
