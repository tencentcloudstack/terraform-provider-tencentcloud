/*
Use this data source to query detailed information of teo botPortraitRules

Example Usage

```hcl
data "tencentcloud_teo_bot_portrait_rules" "botPortraitRules" {
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

func dataSourceTencentCloudTeoBotPortraitRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoBotPortraitRulesRead,
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

			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Portrait rules list.",
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
							Description: "Description of the rule. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"rule_type_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the rule. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"classification_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Classification of the rule. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the rule. Note: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudTeoBotPortraitRulesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_teo_bot_portrait_rules.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var zoneId string
	var entity string

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

	var rules []*teo.PortraitManagedRuleDetail
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := teoService.DescribeTeoBotPortraitRulesByFilter(ctx, paramMap)
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
			if rule.ClassificationId != nil {
				ruleMap["classification_id"] = rule.ClassificationId
			}
			if rule.Status != nil {
				ruleMap["status"] = rule.Status
			}

			ruleList = append(ruleList, ruleMap)
		}
		_ = d.Set("rules", ruleList)
	}

	d.SetId(zoneId + FILED_SP + entity)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), ruleList); e != nil {
			return e
		}
	}
	return nil
}
