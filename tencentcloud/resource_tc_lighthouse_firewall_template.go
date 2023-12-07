package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudLighthouseFirewallTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseFirewallTemplateCreate,
		Read:   resourceTencentCloudLighthouseFirewallTemplateRead,
		Update: resourceTencentCloudLighthouseFirewallTemplateUpdate,
		Delete: resourceTencentCloudLighthouseFirewallTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"template_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Template name.",
			},

			"template_rules": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "List of firewall rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Protocol. Values: TCP, UDP, ICMP, ALL.",
						},
						"port": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Port. Values: ALL, Separate ports, comma-separated discrete ports, minus sign-separated port ranges.",
						},
						"cidr_block": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Network segment or IP (mutually exclusive). The default is `0.0.0.0`, indicating all sources.",
						},
						"action": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Action. Values: ACCEPT, DROP. The default is `ACCEPT`.",
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

func resourceTencentCloudLighthouseFirewallTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_firewall_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = lighthouse.NewCreateFirewallTemplateRequest()
		response   = lighthouse.NewCreateFirewallTemplateResponse()
		templateId string
	)
	if v, ok := d.GetOk("template_name"); ok {
		request.TemplateName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("template_rules"); ok {
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
			request.TemplateRules = append(request.TemplateRules, &firewallRule)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().CreateFirewallTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create lighthouse firewallTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(templateId)

	return resourceTencentCloudLighthouseFirewallTemplateRead(d, meta)
}

func resourceTencentCloudLighthouseFirewallTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_firewall_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	firewallTemplateId := d.Id()

	firewallTemplate, err := service.DescribeFirewallTemplateById(ctx, firewallTemplateId)
	if err != nil {
		return err
	}

	if firewallTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LighthouseFirewallTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if firewallTemplate.TemplateName != nil {
		_ = d.Set("template_name", firewallTemplate.TemplateName)
	}
	templateRules, err := service.DescribeFirewallTemplateRulesById(ctx, firewallTemplateId)
	if err != nil {
		return err
	}
	templateRulesList := []interface{}{}
	for _, templateRule := range templateRules {
		templateRulesMap := map[string]interface{}{}

		if templateRule.FirewallRuleInfo.Protocol != nil {
			templateRulesMap["protocol"] = templateRule.FirewallRuleInfo.Protocol
		}

		if templateRule.FirewallRuleInfo.Port != nil {
			templateRulesMap["port"] = templateRule.FirewallRuleInfo.Port
		}

		if templateRule.FirewallRuleInfo.CidrBlock != nil {
			templateRulesMap["cidr_block"] = templateRule.FirewallRuleInfo.CidrBlock
		}

		if templateRule.FirewallRuleInfo.Action != nil {
			templateRulesMap["action"] = templateRule.FirewallRuleInfo.Action
		}

		if templateRule.FirewallRuleInfo.FirewallRuleDescription != nil {
			templateRulesMap["firewall_rule_description"] = templateRule.FirewallRuleInfo.FirewallRuleDescription
		}

		templateRulesList = append(templateRulesList, templateRulesMap)
	}

	_ = d.Set("template_rules", templateRulesList)

	return nil
}

func resourceTencentCloudLighthouseFirewallTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_firewall_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := lighthouse.NewModifyFirewallTemplateRequest()

	firewallTemplateId := d.Id()

	request.TemplateId = &firewallTemplateId

	immutableArgs := []string{"template_rules"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("template_name") {
		if v, ok := d.GetOk("template_name"); ok {
			request.TemplateName = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().ModifyFirewallTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update lighthouse firewallTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudLighthouseFirewallTemplateRead(d, meta)
}

func resourceTencentCloudLighthouseFirewallTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_firewall_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}
	firewallTemplateId := d.Id()

	if err := service.DeleteFirewallTemplateById(ctx, firewallTemplateId); err != nil {
		return err
	}

	return nil
}
