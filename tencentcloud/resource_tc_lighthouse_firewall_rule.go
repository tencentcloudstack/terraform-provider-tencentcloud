package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
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
		},
	}
}

func resourceTencentCloudLighthouseFirewallRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_firewall_rule.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = lighthouse.NewCreateFirewallRulesRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().CreateFirewallRules(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create lighthouse firewallRule failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudLighthouseFirewallRuleRead(d, meta)
}

func resourceTencentCloudLighthouseFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_firewall_rule.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	firewallRules, err := service.DescribeLighthouseFirewallRuleById(ctx, d.Id())
	if err != nil {
		return err
	}

	if len(firewallRules) == 0 {
		d.SetId("")
		log.Printf("[WARN]%s resource `LighthouseFirewallRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", d.Id())
	firewallRulesList := []interface{}{}
	for _, firewallRule := range firewallRules {

		firewallRulesMap := map[string]interface{}{}

		if firewallRule.Protocol != nil {
			firewallRulesMap["protocol"] = firewallRule.Protocol
		}

		if firewallRule.Port != nil {
			firewallRulesMap["port"] = firewallRule.Port
		}

		if firewallRule.CidrBlock != nil {
			firewallRulesMap["cidr_block"] = firewallRule.CidrBlock
		}

		if firewallRule.Action != nil {
			firewallRulesMap["action"] = firewallRule.Action
		}

		if firewallRule.FirewallRuleDescription != nil {
			firewallRulesMap["firewall_rule_description"] = firewallRule.FirewallRuleDescription
		}

		firewallRulesList = append(firewallRulesList, firewallRulesMap)
	}

	_ = d.Set("firewall_rules", firewallRulesList)

	return nil
}

func resourceTencentCloudLighthouseFirewallRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_firewall_rule.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := lighthouse.NewModifyFirewallRulesRequest()
	hasChanged := false

	if d.HasChange("firewall_rules") {
		if v, ok := d.GetOk("firewall_rules"); ok {
			for _, item := range v.([]interface{}) {
				firewallRule := lighthouse.FirewallRule{}
				dMap := item.(map[string]interface{})
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
		hasChanged = true
	}
	if hasChanged {
		request.InstanceId = helper.String(d.Id())
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
	}

	return resourceTencentCloudLighthouseFirewallRuleRead(d, meta)
}

func resourceTencentCloudLighthouseFirewallRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_firewall_rule.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	firewallRuleInfos, err := service.DescribeLighthouseFirewallRuleById(ctx, d.Id())
	if err != nil {
		return err
	}

	if len(firewallRuleInfos) == 0 {
		d.SetId("")
		log.Printf("[WARN]%s resource `LighthouseFirewallRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	firewallRuleList := make([]*lighthouse.FirewallRule, 0)
	for _, firewallRuleInfo := range firewallRuleInfos {

		var firewallRule lighthouse.FirewallRule

		if firewallRuleInfo.Protocol != nil {
			firewallRule.Protocol = firewallRuleInfo.Protocol
		}

		if firewallRuleInfo.Port != nil {
			firewallRule.Port = firewallRuleInfo.Port
		}

		if firewallRuleInfo.CidrBlock != nil {
			firewallRule.CidrBlock = firewallRuleInfo.CidrBlock
		}

		if firewallRuleInfo.Action != nil {
			firewallRule.Action = firewallRuleInfo.Action
		}

		if firewallRuleInfo.FirewallRuleDescription != nil {
			firewallRule.FirewallRuleDescription = firewallRuleInfo.FirewallRuleDescription
		}

		firewallRuleList = append(firewallRuleList, &firewallRule)
	}

	var (
		request = lighthouse.NewDeleteFirewallRulesRequest()
	)
	request.InstanceId = helper.String(d.Id())
	request.FirewallRules = firewallRuleList

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().DeleteFirewallRules(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s deleted lighthouse firewallRule failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
