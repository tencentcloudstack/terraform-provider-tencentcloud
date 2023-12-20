package cfw

import (
	"context"
	"fmt"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfw "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCfwNatPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfwNatPolicyCreate,
		Read:   resourceTencentCloudCfwNatPolicyRead,
		Update: resourceTencentCloudCfwNatPolicyUpdate,
		Delete: resourceTencentCloudCfwNatPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"source_content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Access source example: net:IP/CIDR(192.168.0.2).",
			},
			"source_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Access source type: for inbound rules, the type can be net, location, vendor, template; for outbound rules, it can be net, instance, tag, template, group.",
			},
			"target_content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Example of access purpose: net: IP/CIDR(192.168.0.2) domain: domain name rules, such as *.qq.com.",
			},
			"target_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Access purpose type: For inbound rules, the type can be net, instance, tag, template, group; for outbound rules, it can be net, location, vendor, template.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Protocol. If Direction=1, optional values: TCP, UDP, ANY; If Direction=0, optional values: TCP, UDP, ICMP, ANY, HTTP, HTTPS, HTTP/HTTPS, SMTP, SMTPS, SMTP/SMTPS, FTP, and DNS.",
			},
			"rule_action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(POLICY_RULE_ACTION),
				Description:  "How the traffic set in the access control policy passes through the cloud firewall. Values: accept: allow; drop: reject; log: observe.",
			},
			"port": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The port for the access control policy. Value: -1/-1: All ports 80: Port 80.",
			},
			"direction": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Rule direction: 1, inbound; 0, outbound.",
			},
			"uuid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique id corresponding to the rule, no need to fill in when creating the rule.",
			},
			"enable": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      POLICY_ENABLE_TRUE,
				ValidateFunc: tccommon.ValidateAllowedStringValue(POLICY_ENABLE),
				Description:  "Rule status, true means enabled, false means disabled. Default is true.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description.",
			},
		},
	}
}

func resourceTencentCloudCfwNatPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_nat_policy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		request           = cfw.NewAddNatAcRuleRequest()
		response          = cfw.NewAddNatAcRuleResponse()
		createNatRuleItem = cfw.CreateNatRuleItem{}
		uuid              string
	)

	if v, ok := d.GetOk("source_content"); ok {
		createNatRuleItem.SourceContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_type"); ok {
		createNatRuleItem.SourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_content"); ok {
		createNatRuleItem.TargetContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_type"); ok {
		createNatRuleItem.TargetType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protocol"); ok {
		createNatRuleItem.Protocol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_action"); ok {
		createNatRuleItem.RuleAction = helper.String(v.(string))
	}

	if v, ok := d.GetOk("port"); ok {
		createNatRuleItem.Port = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("direction"); ok {
		createNatRuleItem.Direction = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("enable"); ok {
		createNatRuleItem.Enable = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		createNatRuleItem.Description = helper.String(v.(string))
	}

	request.Rules = append(request.Rules, &createNatRuleItem)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwClient().AddNatAcRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cfw natPolicy failed, reason:%+v", logId, err)
		return err
	}

	ruleUuid := *response.Response.RuleUuid[0]
	uuid = strconv.FormatInt(ruleUuid, 10)
	d.SetId(uuid)

	return resourceTencentCloudCfwNatPolicyRead(d, meta)
}

func resourceTencentCloudCfwNatPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_nat_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service  = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		ruleUuid = d.Id()
	)

	natPolicy, err := service.DescribeCfwNatPolicyById(ctx, ruleUuid)
	if err != nil {
		return err
	}

	if natPolicy == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CfwNatPolicy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if natPolicy.SourceContent != nil {
		_ = d.Set("source_content", natPolicy.SourceContent)
	}

	if natPolicy.SourceType != nil {
		_ = d.Set("source_type", natPolicy.SourceType)
	}

	if natPolicy.TargetContent != nil {
		_ = d.Set("target_content", natPolicy.TargetContent)
	}

	if natPolicy.TargetType != nil {
		_ = d.Set("target_type", natPolicy.TargetType)
	}

	if natPolicy.Protocol != nil {
		_ = d.Set("protocol", natPolicy.Protocol)
	}

	if natPolicy.Port != nil {
		_ = d.Set("port", natPolicy.Port)
	}

	if natPolicy.Direction != nil {
		_ = d.Set("direction", natPolicy.Direction)
	}

	if natPolicy.Uuid != nil {
		_ = d.Set("uuid", natPolicy.Uuid)
	}

	if natPolicy.Enable != nil {
		_ = d.Set("enable", natPolicy.Enable)
	}

	if natPolicy.Description != nil {
		_ = d.Set("description", natPolicy.Description)
	}

	if natPolicy.Scope != nil {
		_ = d.Set("scope", natPolicy.Scope)
	}

	return nil
}

func resourceTencentCloudCfwNatPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_nat_policy.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		request        = cfw.NewModifyNatAcRuleRequest()
		modifyRuleItem = cfw.CreateNatRuleItem{}
		uuid           = d.Id()
	)

	immutableArgs := []string{"uuid", "direction"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	uuidInt, _ := strconv.ParseInt(uuid, 10, 64)
	modifyRuleItem.Uuid = &uuidInt

	if v, ok := d.GetOk("source_content"); ok {
		modifyRuleItem.SourceContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_type"); ok {
		modifyRuleItem.SourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_content"); ok {
		modifyRuleItem.TargetContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_type"); ok {
		modifyRuleItem.TargetType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protocol"); ok {
		modifyRuleItem.Protocol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_action"); ok {
		modifyRuleItem.RuleAction = helper.String(v.(string))
	}

	if v, ok := d.GetOk("port"); ok {
		modifyRuleItem.Port = helper.String(v.(string))
	}

	if v, ok := d.GetOk("direction"); ok {
		modifyRuleItem.Direction = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("enable"); ok {
		modifyRuleItem.Enable = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		modifyRuleItem.Description = helper.String(v.(string))
	}

	request.Rules = append(request.Rules, &modifyRuleItem)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwClient().ModifyNatAcRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update cfw natPolicy failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCfwNatPolicyRead(d, meta)
}

func resourceTencentCloudCfwNatPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_nat_policy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		uuid    = d.Id()
	)

	if err := service.DeleteCfwNatPolicyById(ctx, uuid); err != nil {
		return err
	}

	return nil
}
