package fwm

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudFwmEdgeAclRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudFwmEdgeAclRuleCreate,
		Read:   resourceTencentCloudFwmEdgeAclRuleRead,
		Update: resourceTencentCloudFwmEdgeAclRuleUpdate,
		Delete: resourceTencentCloudFwmEdgeAclRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Rule group ID.",
			},

			"rules": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Rule list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"order_index": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Rule execution order. The smaller the number, the higher the priority. When creating a rule group, it must start from 1 and strictly increase.",
						},
						"direction": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Rule direction: 0-outbound, 1-inbound.",
						},
						"source_content": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Source address content.",
						},
						"source_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Source address type: ip, domain, template, instance, tag.",
						},
						"target_content": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Target address content.",
						},
						"target_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Target address type: ip, domain, template, instance, tag.",
						},
						"port": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Port, supports single port, port range and comma-separated multi-port.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Protocol type: TCP, UDP, ICMP, ANY.",
						},
						"rule_action": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rule action: accept-pass, drop-block, log-observe.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rule description, 0-256 characters.",
						},
						"scope": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Effective scope: serial, side, all.",
						},
						"belong_member_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Member account ID to which the rule belongs.",
						},
						"param_template_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Parameter template ID.",
						},
					},
				},
			},

			// computed
			"rule_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Rule ID list returned by the cloud API.",
			},
		},
	}
}

func resourceTencentCloudFwmEdgeAclRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_edge_acl_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = fwmv20250611.NewCreateEdgeAclRuleRequest()
		response = fwmv20250611.NewCreateEdgeAclRuleResponse()
		groupId  string
	)

	if v, ok := d.GetOk("group_id"); ok {
		request.GroupId = helper.String(v.(string))
		groupId = v.(string)
	}

	if v, ok := d.GetOk("rules"); ok {
		for _, item := range v.([]interface{}) {
			ruleMap := item.(map[string]interface{})
			rule := fwmv20250611.EdgeAclRuleInfo{}
			if v, ok := ruleMap["order_index"].(int); ok && v != 0 {
				rule.OrderIndex = helper.Int64(int64(v))
			}

			if v, ok := ruleMap["direction"].(int); ok {
				rule.Direction = helper.Int64(int64(v))
			}

			if v, ok := ruleMap["source_content"].(string); ok && v != "" {
				rule.SourceContent = helper.String(v)
			}

			if v, ok := ruleMap["source_type"].(string); ok && v != "" {
				rule.SourceType = helper.String(v)
			}

			if v, ok := ruleMap["target_content"].(string); ok && v != "" {
				rule.TargetContent = helper.String(v)
			}

			if v, ok := ruleMap["target_type"].(string); ok && v != "" {
				rule.TargetType = helper.String(v)
			}

			if v, ok := ruleMap["port"].(string); ok && v != "" {
				rule.Port = helper.String(v)
			}

			if v, ok := ruleMap["protocol"].(string); ok && v != "" {
				rule.Protocol = helper.String(v)
			}

			if v, ok := ruleMap["rule_action"].(string); ok && v != "" {
				rule.RuleAction = helper.String(v)
			}

			if v, ok := ruleMap["description"].(string); ok && v != "" {
				rule.Description = helper.String(v)
			}

			if v, ok := ruleMap["scope"].(string); ok && v != "" {
				rule.Scope = helper.String(v)
			}

			if v, ok := ruleMap["belong_member_id"].(string); ok && v != "" {
				rule.BelongMemberId = helper.String(v)
			}

			if v, ok := ruleMap["param_template_id"].(string); ok && v != "" {
				rule.ParamTemplateId = helper.String(v)
			}

			request.Rules = append(request.Rules, &rule)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().CreateEdgeAclRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create fwm_edge_acl_rule failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create fwm_edge_acl_rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	ruleIds := response.Response.RuleIds
	log.Printf("[DEBUG]%s fwm_edge_acl_rule groupId=%s, d.Id()=%s", logId, groupId, d.Id())
	if len(ruleIds) == 0 {
		return fmt.Errorf("RuleIds is empty.")
	}

	d.SetId(strings.Join([]string{groupId, *ruleIds[0]}, tccommon.FILED_SP))
	_ = d.Set("rule_ids", helper.StringsInterfaces(ruleIds))
	return resourceTencentCloudFwmEdgeAclRuleRead(d, meta)
}

func resourceTencentCloudFwmEdgeAclRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_edge_acl_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = FwmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	groupId := idSplit[0]

	rules, err := service.DescribeEdgeAclRulesByGroupId(ctx, groupId)
	if err != nil {
		return err
	}

	if rules == nil || len(rules) == 0 {
		log.Printf("[CRUD] fwm_edge_acl_rule id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if len(rules) > 0 {
		ruleIds := make([]interface{}, 0, len(rules))
		for _, rule := range rules {
			if rule.RuleId != nil {
				ruleIds = append(ruleIds, *rule.RuleId)
			}
		}
		_ = d.Set("rule_ids", ruleIds)
	}

	rulesList := make([]map[string]interface{}, 0, len(rules))
	for _, rule := range rules {
		ruleMap := map[string]interface{}{}
		if rule.Sequence != nil {
			ruleMap["order_index"] = rule.Sequence
		}

		if rule.Direction != nil {
			ruleMap["direction"] = rule.Direction
		}

		if rule.SourceContent != nil {
			ruleMap["source_content"] = rule.SourceContent
		}

		if rule.SourceType != nil {
			ruleMap["source_type"] = rule.SourceType
		}

		if rule.TargetContent != nil {
			ruleMap["target_content"] = rule.TargetContent
		}

		if rule.TargetType != nil {
			ruleMap["target_type"] = rule.TargetType
		}

		if rule.Port != nil {
			ruleMap["port"] = rule.Port
		}

		if rule.Protocol != nil {
			ruleMap["protocol"] = rule.Protocol
		}

		if rule.RuleAction != nil {
			ruleMap["rule_action"] = rule.RuleAction
		}

		if rule.Description != nil {
			ruleMap["description"] = rule.Description
		}

		if rule.Scope != nil {
			ruleMap["scope"] = rule.Scope
		}

		if rule.ParamTemplateId != nil {
			ruleMap["param_template_id"] = rule.ParamTemplateId
		}

		rulesList = append(rulesList, ruleMap)
	}
	_ = d.Set("rules", rulesList)

	if groupId != "" {
		_ = d.Set("group_id", groupId)
	}

	return nil
}

func resourceTencentCloudFwmEdgeAclRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_edge_acl_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		groupId string
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	groupId = idSplit[0]

	if d.HasChange("rules") {
		oldRuleIdsIf, _ := d.GetChange("rule_ids")
		oldRuleIds := make([]*string, 0, len(oldRuleIdsIf.([]interface{})))
		for _, item := range oldRuleIdsIf.([]interface{}) {
			val := item.(string)
			oldRuleIds = append(oldRuleIds, &val)
		}

		if len(oldRuleIds) > 0 {
			deleteRequest := fwmv20250611.NewDeleteEdgeAclRuleRequest()
			deleteRequest.GroupId = helper.String(groupId)
			deleteRequest.RuleIds = oldRuleIds
			deleteErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().DeleteEdgeAclRuleWithContext(ctx, deleteRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, deleteRequest.GetAction(), deleteRequest.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil {
					return resource.NonRetryableError(fmt.Errorf("Delete fwm_edge_acl_rule failed, Response is nil."))
				}

				return nil
			})

			if deleteErr != nil {
				log.Printf("[CRITAL]%s delete fwm_edge_acl_rule failed when update, reason:%+v", logId, deleteErr)
				return deleteErr
			}
		}

		createRequest := fwmv20250611.NewCreateEdgeAclRuleRequest()
		createRequest.GroupId = helper.String(groupId)
		if v, ok := d.GetOk("rules"); ok {
			for _, item := range v.([]interface{}) {
				ruleMap := item.(map[string]interface{})
				rule := fwmv20250611.EdgeAclRuleInfo{}
				if v, ok := ruleMap["order_index"].(int); ok && v != 0 {
					rule.OrderIndex = helper.Int64(int64(v))
				}

				if v, ok := ruleMap["direction"].(int); ok {
					rule.Direction = helper.Int64(int64(v))
				}

				if v, ok := ruleMap["source_content"].(string); ok && v != "" {
					rule.SourceContent = helper.String(v)
				}

				if v, ok := ruleMap["source_type"].(string); ok && v != "" {
					rule.SourceType = helper.String(v)
				}

				if v, ok := ruleMap["target_content"].(string); ok && v != "" {
					rule.TargetContent = helper.String(v)
				}

				if v, ok := ruleMap["target_type"].(string); ok && v != "" {
					rule.TargetType = helper.String(v)
				}

				if v, ok := ruleMap["port"].(string); ok && v != "" {
					rule.Port = helper.String(v)
				}

				if v, ok := ruleMap["protocol"].(string); ok && v != "" {
					rule.Protocol = helper.String(v)
				}

				if v, ok := ruleMap["rule_action"].(string); ok && v != "" {
					rule.RuleAction = helper.String(v)
				}

				if v, ok := ruleMap["description"].(string); ok && v != "" {
					rule.Description = helper.String(v)
				}

				if v, ok := ruleMap["scope"].(string); ok && v != "" {
					rule.Scope = helper.String(v)
				}

				if v, ok := ruleMap["belong_member_id"].(string); ok && v != "" {
					rule.BelongMemberId = helper.String(v)
				}

				if v, ok := ruleMap["param_template_id"].(string); ok && v != "" {
					rule.ParamTemplateId = helper.String(v)
				}

				createRequest.Rules = append(createRequest.Rules, &rule)
			}
		}

		createResponse := fwmv20250611.NewCreateEdgeAclRuleResponse()
		createErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().CreateEdgeAclRuleWithContext(ctx, createRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, createRequest.GetAction(), createRequest.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Create fwm_edge_acl_rule failed, Response is nil."))
			}

			createResponse = result
			return nil
		})

		if createErr != nil {
			log.Printf("[CRITAL]%s create fwm_edge_acl_rule failed when update, reason:%+v", logId, createErr)
			return createErr
		}

		newRuleIds := createResponse.Response.RuleIds
		log.Printf("[DEBUG]%s fwm_edge_acl_rule update groupId=%s, d.Id()=%s", logId, groupId, d.Id())
		if len(newRuleIds) == 0 {
			return fmt.Errorf("RuleIds is empty after update.")
		}

		d.SetId(strings.Join([]string{groupId, *newRuleIds[0]}, tccommon.FILED_SP))
		_ = d.Set("rule_ids", helper.StringsInterfaces(newRuleIds))
	}

	return resourceTencentCloudFwmEdgeAclRuleRead(d, meta)
}

func resourceTencentCloudFwmEdgeAclRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_edge_acl_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = fwmv20250611.NewDeleteEdgeAclRuleRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	groupId := idSplit[0]
	ruleIdsIf, ok := d.GetOk("rule_ids")
	if !ok {
		return nil
	}

	ruleIds := make([]*string, 0, len(ruleIdsIf.([]interface{})))
	for _, item := range ruleIdsIf.([]interface{}) {
		val := item.(string)
		ruleIds = append(ruleIds, &val)
	}

	if len(ruleIds) == 0 {
		return nil
	}

	request.GroupId = helper.String(groupId)
	request.RuleIds = ruleIds
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().DeleteEdgeAclRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete fwm_edge_acl_rule failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete fwm_edge_acl_rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
