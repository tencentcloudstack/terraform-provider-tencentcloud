package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudGaapSecurityRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapSecurityRulesRead,
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"action": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"ACCEPT", "DROP"}),
			},
			"cidr_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCidrIp,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value == "ALL" {
						return
					}
					match, _ := regexp.MatchString("^(\\d{1,5},)*\\d{1,5}$|^\\d{1,5}-\\d{1,5}$", value)
					if !match {
						errors = append(errors, fmt.Errorf("%s example: 53、80,443、80-90, Not configured to represent all ports", k))
					}
					return
				},
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"ALL", "TCP", "UDP"}),
			},

			// computed
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action": {
							Type:     schema.TypeString,
							Computed: true,
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
		ruleId = stringToPointer(raw.(string))
	}
	if raw, ok := d.GetOk("action"); ok {
		action = stringToPointer(raw.(string))
	}
	if raw, ok := d.GetOk("cidr_ip"); ok {
		cidrIp = stringToPointer(raw.(string))
	}
	if raw, ok := d.GetOk("name"); ok {
		name = stringToPointer(raw.(string))
	}
	if raw, ok := d.GetOk("port"); ok {
		port = stringToPointer(raw.(string))
	}
	if raw, ok := d.GetOk("protocol"); ok {
		protocol = stringToPointer(raw.(string))
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

	d.Set("rules", rules)
	d.SetId(dataResourceIdsHash(ids))

	return nil
}
