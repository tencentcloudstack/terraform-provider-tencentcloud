package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfw "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCfwEdgePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfwEdgePolicyCreate,
		Read:   resourceTencentCloudCfwEdgePolicyRead,
		Update: resourceTencentCloudCfwEdgePolicyUpdate,
		Delete: resourceTencentCloudCfwEdgePolicyDelete,
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
				Description: "Protocol. If Direction=1 && Scope=serial, optional values: TCP UDP ICMP ANY HTTP HTTPS HTTP/HTTPS SMTP SMTPS SMTP/SMTPS FTP DNS; If Direction=1 && Scope!=serial, optional values: TCP; If Direction=0 && Scope=serial, optional values: TCP UDP ICMP ANY HTTP HTTPS HTTP/HTTPS SMTP SMTPS SMTP/SMTPS FTP DNS; If Direction=0 && Scope!=serial, optional values: TCP HTTP/HTTPS TLS/SSL.",
			},
			"rule_action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(POLICY_RULE_ACTION),
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
				ValidateFunc: validateAllowedStringValue(POLICY_ENABLE),
				Description:  "Rule status, true means enabled, false means disabled. Default is true.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description.",
			},
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      POLICY_SCOPE_ALL,
				ValidateFunc: validateAllowedStringValue(POLICY_SCOPE),
				Description:  "Effective range. serial: serial; side: bypass; all: global, Default is all.",
			},
		},
	}
}

func resourceTencentCloudCfwEdgePolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_edge_policy.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId          = getLogId(contextNil)
		request        = cfw.NewAddAclRuleRequest()
		response       = cfw.NewAddAclRuleResponse()
		createRuleItem = cfw.CreateRuleItem{}
		uuid           string
	)

	if v, ok := d.GetOk("source_content"); ok {
		createRuleItem.SourceContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_type"); ok {
		createRuleItem.SourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_content"); ok {
		createRuleItem.TargetContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_type"); ok {
		createRuleItem.TargetType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protocol"); ok {
		createRuleItem.Protocol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_action"); ok {
		createRuleItem.RuleAction = helper.String(v.(string))
	}

	if v, ok := d.GetOk("port"); ok {
		createRuleItem.Port = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("direction"); ok {
		createRuleItem.Direction = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("enable"); ok {
		createRuleItem.Enable = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		createRuleItem.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("scope"); ok {
		createRuleItem.Scope = helper.String(v.(string))
	}

	request.Rules = append(request.Rules, &createRuleItem)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCfwClient().AddAclRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cfw edgePolicy failed, reason:%+v", logId, err)
		return err
	}

	ruleUuid := *response.Response.RuleUuid[0]
	uuid = strconv.FormatInt(ruleUuid, 10)
	d.SetId(uuid)

	return resourceTencentCloudCfwEdgePolicyRead(d, meta)
}

func resourceTencentCloudCfwEdgePolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_edge_policy.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = CfwService{client: meta.(*TencentCloudClient).apiV3Conn}
		ruleUuid   = d.Id()
		sourceType string
		targetType string
	)

	edgePolicy, err := service.DescribeCfwEdgePolicyById(ctx, ruleUuid)
	if err != nil {
		return err
	}

	if edgePolicy == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CfwEdgePolicy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if edgePolicy.SourceType != nil {
		_ = d.Set("source_type", edgePolicy.SourceType)
		sourceType = *edgePolicy.SourceType
	}

	if edgePolicy.SourceContent != nil {
		if sourceType == "tag" {
			params := strings.Split(*edgePolicy.SourceContent, "|")
			key := params[0]
			value := params[1]
			var obj SourceContentJson
			obj.Key = key
			obj.Value = value
			tmpStr, _ := json.Marshal(obj)
			_ = d.Set("source_content", string(tmpStr))
		} else {
			_ = d.Set("source_content", edgePolicy.SourceContent)
		}
	}

	if edgePolicy.TargetType != nil {
		_ = d.Set("target_type", edgePolicy.TargetType)
		targetType = *edgePolicy.TargetType
	}

	if edgePolicy.TargetContent != nil {
		if targetType == "tag" {
			params := strings.Split(*edgePolicy.TargetContent, "|")
			key := params[0]
			value := params[1]
			var obj TargetContentJson
			obj.Key = key
			obj.Value = value
			tmpStr, _ := json.Marshal(obj)
			_ = d.Set("target_content", string(tmpStr))
		} else {
			_ = d.Set("target_content", edgePolicy.TargetContent)
		}
	}

	if edgePolicy.Protocol != nil {
		_ = d.Set("protocol", edgePolicy.Protocol)
	}

	if edgePolicy.Port != nil {
		_ = d.Set("port", edgePolicy.Port)
	}

	if edgePolicy.Direction != nil {
		_ = d.Set("direction", edgePolicy.Direction)
	}

	if edgePolicy.Uuid != nil {
		_ = d.Set("uuid", edgePolicy.Uuid)
	}

	if edgePolicy.Enable != nil {
		_ = d.Set("enable", edgePolicy.Enable)
	}

	if edgePolicy.Description != nil {
		_ = d.Set("description", edgePolicy.Description)
	}

	if edgePolicy.Scope != nil {
		_ = d.Set("scope", edgePolicy.Scope)
	}

	return nil
}

func resourceTencentCloudCfwEdgePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_edge_policy.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId          = getLogId(contextNil)
		request        = cfw.NewModifyAclRuleRequest()
		modifyRuleItem = cfw.CreateRuleItem{}
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

	if v, ok := d.GetOk("scope"); ok {
		modifyRuleItem.Scope = helper.String(v.(string))
	}

	request.Rules = append(request.Rules, &modifyRuleItem)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCfwClient().ModifyAclRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update cfw edgePolicy failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCfwEdgePolicyRead(d, meta)
}

func resourceTencentCloudCfwEdgePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_edge_policy.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = CfwService{client: meta.(*TencentCloudClient).apiV3Conn}
		uuid    = d.Id()
	)

	if err := service.DeleteCfwEdgePolicyById(ctx, uuid); err != nil {
		return err
	}

	return nil
}
