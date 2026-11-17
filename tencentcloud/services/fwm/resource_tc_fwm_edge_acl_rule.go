package fwm

import (
	"context"
	"fmt"
	"log"
	"strconv"
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
				Description: "Edge ACL rule group ID. The rule group is the parent container and must already exist before creating the rule.",
			},

			"direction": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Rule direction: 0-egress, 1-ingress. Changing this will recreate the rule.",
			},

			"source_content": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Source address content. Format depends on source_type: IP/CIDR for ip, domain for domain, template ID for template, instance ID list (comma separated) for instance, tag key:value for tag.",
			},

			"source_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Source address type: ip, domain, template, instance, tag. Changing this will recreate the rule.",
			},

			"target_content": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Target address content. Format is the same as source_content.",
			},

			"target_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Target address type: ip, domain, template, instance, tag. Changing this will recreate the rule.",
			},

			"port": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Port. Supports single port, port range and comma separated multi ports, such as 80, 1-65535, 80,443,8080.",
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
				Description: "Rule description, length 0-256 characters.",
			},

			"scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Effective scope: serial, side, all.",
			},

			"belong_member_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Member account ID that the rule belongs to (used in multi-account scenarios). Required when source_type or target_type is instance or tag.",
			},

			"param_template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Parameter template ID.",
			},

			// computed
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule ID.",
			},

			"sequence": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Rule execution sequence returned by the server.",
			},
		},
	}
}

func resourceTencentCloudFwmEdgeAclRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_edge_acl_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = fwmv20250611.NewCreateEdgeAclRuleRequest()
		response  = fwmv20250611.NewCreateEdgeAclRuleResponse()
		groupId   string
		ruleId    string
		direction int64
	)

	if v, ok := d.GetOk("group_id"); ok {
		request.GroupId = helper.String(v.(string))
		groupId = v.(string)
	}

	rule := fwmv20250611.EdgeAclRuleInfo{}
	if v, ok := d.GetOkExists("direction"); ok {
		direction = int64(v.(int))
		rule.Direction = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("source_content"); ok {
		rule.SourceContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_type"); ok {
		rule.SourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_content"); ok {
		rule.TargetContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_type"); ok {
		rule.TargetType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("port"); ok {
		rule.Port = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protocol"); ok {
		rule.Protocol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_action"); ok {
		rule.RuleAction = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		rule.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("scope"); ok {
		rule.Scope = helper.String(v.(string))
	}

	if v, ok := d.GetOk("belong_member_id"); ok {
		rule.BelongMemberId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("param_template_id"); ok {
		rule.ParamTemplateId = helper.String(v.(string))
	}

	rule.OrderIndex = helper.IntInt64(1)
	request.Rules = []*fwmv20250611.EdgeAclRuleInfo{&rule}
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().CreateEdgeAclRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create fwm edge acl rule failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create fwm edge acl rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s create fwm edge acl rule success, logId=%s, d.Id()=%s", logId, logId, d.Id())

	if len(response.Response.RuleIds) == 0 || *response.Response.RuleIds[0] == "" {
		return fmt.Errorf("Create fwm edge acl rule failed, RuleIds is empty or first element is empty.")
	}

	ruleId = *response.Response.RuleIds[0]
	d.SetId(strings.Join([]string{groupId, ruleId, strconv.FormatInt(direction, 10)}, tccommon.FILED_SP))
	return resourceTencentCloudFwmEdgeAclRuleRead(d, meta)
}

func resourceTencentCloudFwmEdgeAclRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_edge_acl_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = NewFwmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	groupId, ruleId, direction, err := parseFwmEdgeAclRuleId(d.Id())
	if err != nil {
		return err
	}

	respData, err := service.DescribeFwmEdgeAclRuleById(ctx, groupId, ruleId, direction)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_fwm_edge_acl_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("group_id", groupId)

	if respData.Direction != nil {
		_ = d.Set("direction", respData.Direction)
	}

	if respData.SourceContent != nil {
		_ = d.Set("source_content", respData.SourceContent)
	}

	if respData.SourceType != nil {
		_ = d.Set("source_type", respData.SourceType)
	}

	if respData.TargetContent != nil {
		_ = d.Set("target_content", respData.TargetContent)
	}

	if respData.TargetType != nil {
		_ = d.Set("target_type", respData.TargetType)
	}

	if respData.Port != nil {
		_ = d.Set("port", respData.Port)
	}

	if respData.Protocol != nil {
		_ = d.Set("protocol", respData.Protocol)
	}

	if respData.RuleAction != nil {
		_ = d.Set("rule_action", respData.RuleAction)
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.Scope != nil {
		_ = d.Set("scope", respData.Scope)
	}

	if respData.BelongMember != nil && respData.BelongMember.MemberId != nil {
		_ = d.Set("belong_member_id", respData.BelongMember.MemberId)
	}

	if respData.ParamTemplateId != nil {
		_ = d.Set("param_template_id", respData.ParamTemplateId)
	}

	if respData.RuleId != nil {
		_ = d.Set("rule_id", respData.RuleId)
	}

	if respData.Sequence != nil {
		_ = d.Set("sequence", respData.Sequence)
	}

	return nil
}

func resourceTencentCloudFwmEdgeAclRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_edge_acl_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	groupId, ruleId, direction, err := parseFwmEdgeAclRuleId(d.Id())
	if err != nil {
		return err
	}

	needChange := false
	mutableArgs := []string{"source_content", "target_content", "port", "protocol", "rule_action", "description", "scope", "belong_member_id", "param_template_id"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := fwmv20250611.NewModifyEdgeAclRuleRequest()
		rule := fwmv20250611.EdgeAclRuleInfo{}
		if v, ok := d.GetOk("source_content"); ok {
			rule.SourceContent = helper.String(v.(string))
		}

		if v, ok := d.GetOk("source_type"); ok {
			rule.SourceType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("target_content"); ok {
			rule.TargetContent = helper.String(v.(string))
		}

		if v, ok := d.GetOk("target_type"); ok {
			rule.TargetType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("port"); ok {
			rule.Port = helper.String(v.(string))
		}

		if v, ok := d.GetOk("protocol"); ok {
			rule.Protocol = helper.String(v.(string))
		}

		if v, ok := d.GetOk("rule_action"); ok {
			rule.RuleAction = helper.String(v.(string))
		}

		if v, ok := d.GetOk("description"); ok {
			rule.Description = helper.String(v.(string))
		}

		if v, ok := d.GetOk("scope"); ok {
			rule.Scope = helper.String(v.(string))
		}

		if v, ok := d.GetOk("belong_member_id"); ok {
			rule.BelongMemberId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("param_template_id"); ok {
			rule.ParamTemplateId = helper.String(v.(string))
		}

		rule.RuleId = helper.String(ruleId)
		rule.Direction = helper.Int64(direction)
		rule.OrderIndex = helper.IntInt64(1)
		request.Rule = &rule
		request.GroupId = helper.String(groupId)
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().ModifyEdgeAclRuleWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify fwm edge acl rule failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update fwm edge acl rule failed, reason:%+v", logId, reqErr)
			return reqErr
		}
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

	groupId, ruleId, _, err := parseFwmEdgeAclRuleId(d.Id())
	if err != nil {
		return err
	}

	request.GroupId = helper.String(groupId)
	request.RuleIds = []*string{helper.String(ruleId)}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().DeleteEdgeAclRuleWithContext(ctx, request)
		if e != nil {
			if strings.Contains(e.Error(), "ResourceNotFound") {
				log.Printf("[WARN]%s delete fwm edge acl rule, rule not found, treat as success, id=%s", logId, d.Id())
				return nil
			}
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete fwm edge acl rule failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete fwm edge acl rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

// parseFwmEdgeAclRuleId splits the resource id <groupId>#<ruleId>#<direction>.
func parseFwmEdgeAclRuleId(id string) (groupId, ruleId string, direction int64, err error) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 3 {
		return "", "", 0, fmt.Errorf("id is broken, %s", id)
	}

	direction, err = strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return "", "", 0, fmt.Errorf("id is broken, invalid direction in %s: %s", id, err.Error())
	}

	return parts[0], parts[1], direction, nil
}
