/*
Use this data source to query security policy rule.

Example Usage

```hcl
resource "tencentcloud_gaap_proxy" "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource "tencentcloud_gaap_security_policy" "foo" {
  proxy_id = tencentcloud_gaap_proxy.foo.id
  action   = "ACCEPT"
}

resource "tencentcloud_gaap_security_rule" "foo" {
  policy_id = tencentcloud_gaap_security_policy.foo.id
  name      = "ci-test-gaap-s-rule"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}

data "tencentcloud_gaap_security_rules" "protocol" {
  policy_id = tencentcloud_gaap_security_policy.foo.id
  protocol  = tencentcloud_gaap_security_rule.foo.protocol
}
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudGaapSecurityRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapSecurityRulesRead,
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the security policy to be queried.",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the security policy rules to be queried.",
			},
			"action": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"ACCEPT", "DROP"}),
				Description:  "Policy of the rule to be queried.",
			},
			"cidr_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCidrIp,
				Description:  "A network address block of the request source to be queried.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the security policy rule to be queried.",
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value == "ALL" {
						return
					}
					if !regexp.MustCompile(`^(\d{1,5},)*\d{1,5}$|^\d{1,5}-\d{1,5}$`).MatchString(value) {
						errors = append(errors, fmt.Errorf("%s example: `53`, `80,443` and `80-90`, Not configured to represent all ports", k))
					}
					return
				},
				Description: "Port of the security policy rule to be queried.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"ALL", "TCP", "UDP"}),
				Description:  "Protocol of the security policy rule to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// computed
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of security policy rule. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the security policy rule.",
						},
						"cidr_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A network address block of the request source.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the security policy rule.",
						},
						"port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Port of the security policy rule.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protocol of the security policy rule.",
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy of the rule.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudGaapSecurityRulesRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("data_source.tencentcloud_gaap_security_rules.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	policyId := d.Get("policy_id").(string)
	var (
		ruleId   *string
		action   *string
		cidrIp   *string
		name     *string
		port     *string
		protocol *string
		ids      []string
		rules    []map[string]interface{}
	)

	if raw, ok := d.GetOk("rule_id"); ok {
		ruleId = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("action"); ok {
		action = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("cidr_ip"); ok {
		cidrIp = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("name"); ok {
		name = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("port"); ok {
		port = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("protocol"); ok {
		protocol = helper.String(raw.(string))
	}

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	securityRules, err := service.DescribeSecurityRules(ctx, policyId)
	if err != nil {
		return err
	}

	ids = make([]string, 0, len(securityRules))
	rules = make([]map[string]interface{}, 0, len(securityRules))

	for _, rule := range securityRules {
		if rule.RuleId == nil {
			return errors.New("rule id is nil")
		}
		if rule.SourceCidr == nil {
			return errors.New("rule cidr ip is nil")
		}
		if rule.AliasName == nil {
			return errors.New("rule name is nil")
		}
		if rule.DestPortRange == nil {
			return errors.New("rule port is nil")
		}
		if rule.Protocol == nil {
			return errors.New("rule protocol is nil")
		}
		if rule.Action == nil {
			return errors.New("rule action is nil")
		}

		if ruleId != nil && *ruleId != *rule.RuleId {
			continue
		}
		if action != nil && *action != *rule.Action {
			continue
		}
		if cidrIp != nil && *cidrIp != *rule.SourceCidr {
			continue
		}
		if name != nil && *name != *rule.AliasName {
			continue
		}
		if port != nil && *port != *rule.DestPortRange {
			continue
		}
		if protocol != nil && *protocol != *rule.Protocol {
			continue
		}

		ids = append(ids, *rule.RuleId)
		rules = append(rules, map[string]interface{}{
			"id":       *rule.RuleId,
			"cidr_ip":  *rule.SourceCidr,
			"name":     *rule.AliasName,
			"port":     *rule.DestPortRange,
			"protocol": *rule.Protocol,
			"action":   *rule.Action,
		})
	}

	_ = d.Set("rules", rules)
	d.SetId(helper.DataResourceIdsHash(ids))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), rules); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
			return err
		}
	}

	return nil
}
