/*
Use this data source to query dayu L7 rules
Example Usage
```hcl
data "tencentcloud_dayu_l7_rules" "domain_test" {
  resource_type = tencentcloud_dayu_l7_rule.test_rule.resource_type
  resource_id      = tencentcloud_dayu_l7_rule.test_rule.resource_id
  domain		= tencentcloud_dayu_l7_rule.test_rule.domain
}
data "tencentcloud_dayu_l7_rules" "id_test" {
  resource_type = tencentcloud_dayu_l7_rule.test_rule.resource_type
  resource_id      = tencentcloud_dayu_l7_rule.test_rule.resource_id
  rule_id		= tencentcloud_dayu_l7_rule.test_rule.rule_id
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDayuL7Rules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDayuL7RulesRead,
		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(DAYU_RESOURCE_TYPE),
				Description:  "Type of the resource that the L7 rule works for, valid values are `bgpip`, `bgp`, `bgp-multip`, `net`.",
			},
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of the resource that the L7 rule works for.",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of the L7 rule to be queried.",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Domain of the L7 rule to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of L7 rules. Each element contains the following attributes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain that the 7 layer rule works for.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protocol of the rule.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the rule.",
						},
						"switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicate the rule will take effect or not.",
						},
						"source_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Source type, 1 for source of host, 2 for source of ip.",
						},
						"source_list": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type:        schema.TypeString,
								Description: "Source ip or domain.",
							},
							Description: "Source list of the rule.",
						},
						"ssl_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SSL id.",
						},
						"health_check_switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether health check is enabled.",
						},
						"health_check_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Interval time of health check.",
						},
						"health_check_health_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Health threshold of health check.",
						},
						"health_check_unhealth_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unhealthy threshold of health check.",
						},
						"health_check_code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "HTTP Status Code. 1 means the return value '1xx' is health. 2 means the return value '2xx' is health. 4 means the return value '3xx' is health. 8 means the return value '4xx' is health. 16 means the return value '5xx' is health. If you want multiple return codes to indicate health, need to add the corresponding values.",
						},
						"health_check_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path of health check.",
						},
						"health_check_method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Methods of health check.",
						},
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the 7 layer rule.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Status of the rule. 0 for create/modify success, 2 for create/modify fail, 3 for delete success, 5 for waiting to be created/modified, 7 for waiting to be deleted and 8 for waiting to get SSL id.",
						},
						"threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Threshold of the rule.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDayuL7RulesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dayu_l7_rules.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := DayuService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	resourceType := d.Get("resource_type").(string)
	resourceId := d.Get("resource_id").(string)
	ruleId := ""
	if v, ok := d.GetOk("rule_id"); ok {
		ruleId = v.(string)
	}
	domain := ""
	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
	}
	protocol := ""
	if v, ok := d.GetOk("protocol"); ok {
		protocol = v.(string)
	}

	rules, healths, _, err := service.DescribeL7Rules(ctx, resourceType, resourceId, domain, ruleId, protocol)
	if err != nil {
		rules, healths, _, err = service.DescribeL7Rules(ctx, resourceType, resourceId, domain, ruleId, protocol)
	}

	if err != nil {
		return err
	}

	list := make([]map[string]interface{}, 0, len(rules))
	ids := make([]string, 0, len(rules))

	listItem := make(map[string]interface{})
	for k, rule := range rules {
		health := healths[k]
		listItem["name"] = *rule.RuleName
		listItem["domain"] = *rule.Domain
		listItem["ssl_id"] = *rule.SSLId
		listItem["rule_id"] = *rule.RuleId
		listItem["protocol"] = *rule.Protocol
		listItem["source_type"] = int(*rule.SourceType)
		listItem["status"] = int(*rule.Status)
		listItem["threshold"] = int(*rule.CCThreshold)

		if *rule.Protocol == DAYU_L7_RULE_PROTOCOL_HTTPS {
			listItem["switch"] = intToBool(int(*rule.CCEnable))
		} else {
			listItem["switch"] = intToBool(int(*rule.CCStatus))
		}
		sourceList := make([]*string, 0, len(rule.SourceList))
		for _, v := range rule.SourceList {
			sourceList = append(sourceList, v.Source)
		}
		listItem["source_list"] = helper.StringsInterfaces(sourceList)

		if health.Enable != nil {
			listItem["health_check_switch"] = intToBool(int(*health.Enable))
		}
		if health.Url != nil {
			listItem["health_check_path"] = *health.Url
		}
		if health.StatusCode != nil {
			listItem["health_check_code"] = int(*health.StatusCode)
		}
		if health.Interval != nil {
			listItem["health_check_interval"] = int(*health.Interval)
		}
		if health.KickNum != nil {
			listItem["health_check_unhealth_num"] = int(*health.KickNum)
		}
		if health.AliveNum != nil {
			listItem["health_check_health_num"] = int(*health.AliveNum)
		}
		if health.Method != nil {
			listItem["health_check_method"] = *health.Method
		}
		list = append(list, listItem)
		ids = append(ids, listItem["rule_id"].(string))
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("list", list); e != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s\n", logId, e.Error())
		return e
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return writeToFile(output.(string), list)
	}
	return nil

}
