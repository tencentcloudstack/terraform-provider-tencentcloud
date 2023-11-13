/*
Provides a resource to create a lighthouse firewall_rule

Example Usage

```hcl
resource "tencentcloud_lighthouse_firewall_rule" "firewall_rule" {
  instance_id = "lhins-acb1234"
  firewall_rules {
		protocol = "TCP"
		port = "80"
		cidr_block = "22"
		action = "ACCEPT"
		firewall_rule_description = "description"

  }
  firewall_version = 1
}
```

Import

lighthouse firewall_rule can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_firewall_rule.firewall_rule firewall_rule_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudLighthouseFirewallRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseFirewallRuleCreate,
		Read:   resourceTencentCloudLighthouseFirewallRuleRead,
		Update: resourceTencentCloudLighthouseFirewallRuleUpdate,
		Delete: resourceTencentCloudLighthouseFirewallRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID。.",
			},

			"firewall_rules": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Firewall rule list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Protocol. Valid values are TCP, UDP, ICMP, ALL.",
						},
						"port": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Port. Valid values are ALL, one single port, multiple ports separated by commas, or port range indicated by a minus sign.",
						},
						"cidr_block": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IP range or IP (mutually exclusive). Default value is 0.0.0.0/0, which indicates all sources.",
						},
						"action": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Valid values are ACCEPT, DROP. Default value is ACCEPT.",
						},
						"firewall_rule_description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Firewall rule description.",
						},
					},
				},
			},

			"firewall_version": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Current firewall version number. Every time you update the firewall rule version, it will be automatically increased by 1 to prevent the rule from expiring. If it is left empty, conflicts will not be considered.",
			},
		},
	}
}

func resourceTencentCloudLighthouseFirewallRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_firewall_rule.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = lighthouse.NewCreateFirewallRulesRequest()
		response = lighthouse.NewCreateFirewallRulesResponse()
		protocol string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("firewall_rules"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			firewallRule := lighthouse.FirewallRule{}
			if v, ok := dMap["protocol"]; ok {
				firewallRule.Protocol = helper.String(v.(string))
			}
			if v, ok := dMap["port"]; ok {
				firewallRule.Port = helper.String(v.(string))
			}
			if v, ok := dMap["cidr_block"]; ok {
				firewallRule.CidrBlock = helper.String(v.(string))
			}
			if v, ok := dMap["action"]; ok {
				firewallRule.Action = helper.String(v.(string))
			}
			if v, ok := dMap["firewall_rule_description"]; ok {
				firewallRule.FirewallRuleDescription = helper.String(v.(string))
			}
			request.FirewallRules = append(request.FirewallRules, &firewallRule)
		}
	}

	if v, ok := d.GetOkExists("firewall_version"); ok {
		request.FirewallVersion = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().CreateFirewallRules(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create lighthouse firewallRule failed, reason:%+v", logId, err)
		return err
	}

	protocol = *response.Response.Protocol
	d.SetId(protocol)

	return resourceTencentCloudLighthouseFirewallRuleRead(d, meta)
}

func resourceTencentCloudLighthouseFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_firewall_rule.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LighthouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	firewallRuleId := d.Id()

	firewallRule, err := service.DescribeLighthouseFirewallRuleById(ctx, protocol)
	if err != nil {
		return err
	}

	if firewallRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LighthouseFirewallRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if firewallRule.InstanceId != nil {
		_ = d.Set("instance_id", firewallRule.InstanceId)
	}

	if firewallRule.FirewallRules != nil {
		firewallRulesList := []interface{}{}
		for _, firewallRules := range firewallRule.FirewallRules {
			firewallRulesMap := map[string]interface{}{}

			if firewallRule.FirewallRules.Protocol != nil {
				firewallRulesMap["protocol"] = firewallRule.FirewallRules.Protocol
			}

			if firewallRule.FirewallRules.Port != nil {
				firewallRulesMap["port"] = firewallRule.FirewallRules.Port
			}

			if firewallRule.FirewallRules.CidrBlock != nil {
				firewallRulesMap["cidr_block"] = firewallRule.FirewallRules.CidrBlock
			}

			if firewallRule.FirewallRules.Action != nil {
				firewallRulesMap["action"] = firewallRule.FirewallRules.Action
			}

			if firewallRule.FirewallRules.FirewallRuleDescription != nil {
				firewallRulesMap["firewall_rule_description"] = firewallRule.FirewallRules.FirewallRuleDescription
			}

			firewallRulesList = append(firewallRulesList, firewallRulesMap)
		}

		_ = d.Set("firewall_rules", firewallRulesList)

	}

	if firewallRule.FirewallVersion != nil {
		_ = d.Set("firewall_version", firewallRule.FirewallVersion)
	}

	return nil
}

func resourceTencentCloudLighthouseFirewallRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_firewall_rule.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		modifyFirewallRulesRequest  = lighthouse.NewModifyFirewallRulesRequest()
		modifyFirewallRulesResponse = lighthouse.NewModifyFirewallRulesResponse()
	)

	firewallRuleId := d.Id()

	request.Protocol = &protocol

	immutableArgs := []string{"instance_id", "firewall_rules", "firewall_version"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("instance_id") {
		if v, ok := d.GetOk("instance_id"); ok {
			request.InstanceId = helper.String(v.(string))
		}
	}

	if d.HasChange("firewall_rules") {
		if v, ok := d.GetOk("firewall_rules"); ok {
			for _, item := range v.([]interface{}) {
				firewallRule := lighthouse.FirewallRule{}
				if v, ok := dMap["protocol"]; ok {
					firewallRule.Protocol = helper.String(v.(string))
				}
				if v, ok := dMap["port"]; ok {
					firewallRule.Port = helper.String(v.(string))
				}
				if v, ok := dMap["cidr_block"]; ok {
					firewallRule.CidrBlock = helper.String(v.(string))
				}
				if v, ok := dMap["action"]; ok {
					firewallRule.Action = helper.String(v.(string))
				}
				if v, ok := dMap["firewall_rule_description"]; ok {
					firewallRule.FirewallRuleDescription = helper.String(v.(string))
				}
				request.FirewallRules = append(request.FirewallRules, &firewallRule)
			}
		}
	}

	if d.HasChange("firewall_version") {
		if v, ok := d.GetOkExists("firewall_version"); ok {
			request.FirewallVersion = helper.IntUint64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().ModifyFirewallRules(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update lighthouse firewallRule failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudLighthouseFirewallRuleRead(d, meta)
}

func resourceTencentCloudLighthouseFirewallRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_firewall_rule.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LighthouseService{client: meta.(*TencentCloudClient).apiV3Conn}
	firewallRuleId := d.Id()

	if err := service.DeleteLighthouseFirewallRuleById(ctx, protocol); err != nil {
		return err
	}

	return nil
}
