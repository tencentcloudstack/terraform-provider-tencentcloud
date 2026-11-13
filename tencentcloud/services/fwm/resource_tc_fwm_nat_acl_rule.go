package fwm

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudFwmNatAclRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudFwmNatAclRuleCreate,
		Read:   resourceTencentCloudFwmNatAclRuleRead,
		Update: resourceTencentCloudFwmNatAclRuleUpdate,
		Delete: resourceTencentCloudFwmNatAclRuleDelete,
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

			"direction": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Rule direction: 0-outbound, 1-inbound. It is also part of the resource id, so changing it will recreate the rule.",
			},

			"source_content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Source address content.",
			},

			"source_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Source type: ip/net/url/domain/template/addressbook/instance/tag/location/vendor/dnsparse.",
			},

			"target_content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Target address content.",
			},

			"target_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Target type (same enum as source_type).",
			},

			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Protocol: TCP/UDP/ICMP/ANY/HTTP/HTTPS/DNS/FTP...",
			},

			"rule_action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Action: accept/drop/log (lowercase).",
			},

			"scope": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Effective range: ALL/region/nat-instance-id.",
			},

			"port": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Port e.g. 80, 80,443, 80-90; ANY = -1/-1.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule description.",
			},

			"param_template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Port template ID.",
			},

			"belong_member_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Member account ID (required when source/target type in instance/tag/template/addressbook).",
			},

			// computed
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule ID, it is part of the resource id `group_id#rule_id#direction`.",
			},

			"sequence": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Rule priority sequence returned by the server.",
			},
		},
	}
}

func resourceTencentCloudFwmNatAclRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_nat_acl_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = fwmv20250611.NewCreateNatAclRuleRequest()
		response  = fwmv20250611.NewCreateNatAclRuleResponse()
		groupId   string
		ruleId    string
		direction int64
	)

	if v, ok := d.GetOk("group_id"); ok {
		request.GroupId = helper.String(v.(string))
		groupId = v.(string)
	}

	rule := fwmv20250611.NatAclRule{}
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
	if v, ok := d.GetOk("protocol"); ok {
		rule.Protocol = helper.String(v.(string))
	}
	if v, ok := d.GetOk("rule_action"); ok {
		rule.RuleAction = helper.String(v.(string))
	}
	if v, ok := d.GetOk("scope"); ok {
		rule.Scope = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("direction"); ok {
		direction = int64(v.(int))
		rule.Direction = helper.Int64(direction)
	}
	if v, ok := d.GetOk("port"); ok {
		rule.Port = helper.String(v.(string))
	}
	if v, ok := d.GetOk("description"); ok {
		rule.Description = helper.String(v.(string))
	}
	if v, ok := d.GetOk("param_template_id"); ok {
		rule.ParamTemplateId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("belong_member_id"); ok {
		rule.BelongMemberId = helper.String(v.(string))
	}

	rule.OrderIndex = helper.Int64(int64(1))
	request.Rules = []*fwmv20250611.NatAclRule{&rule}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().CreateNatAclRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create fwm nat acl rule failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create fwm nat acl rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s create fwm nat acl rule success, logId=%s, id=%s", logId, logId, d.Id())
	if len(response.Response.RuleIds) < 1 {
		return fmt.Errorf("Create fwm nat acl rule failed, RuleIds is empty.")
	}
	ruleId = *response.Response.RuleIds[0]
	d.SetId(strings.Join([]string{groupId, ruleId, strconv.FormatInt(direction, 10)}, tccommon.FILED_SP))
	return resourceTencentCloudFwmNatAclRuleRead(d, meta)
}

func resourceTencentCloudFwmNatAclRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_nat_acl_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = FwmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	groupId, ruleId, direction, err := parseFwmNatAclRuleId(d.Id())
	if err != nil {
		return err
	}

	respData, err := service.DescribeFwmNatAclRuleById(ctx, groupId, ruleId, direction)
	if err != nil {
		if sdkerr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			if sdkerr.Code == "ResourceNotFound" {
				log.Printf("[CRUD] fwm_nat_acl_rule id=%s", d.Id())
				d.SetId("")
				return nil
			}
		}
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_fwm_nat_acl_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}
	_ = d.Set("group_id", groupId)
	if respData.RuleId != nil {
		_ = d.Set("rule_id", respData.RuleId)
	}
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
	if respData.Protocol != nil {
		_ = d.Set("protocol", respData.Protocol)
	}
	if respData.RuleAction != nil {
		_ = d.Set("rule_action", respData.RuleAction)
	}
	if respData.Sequence != nil {
		_ = d.Set("sequence", respData.Sequence)
	}
	if respData.Scope != nil {
		_ = d.Set("scope", respData.Scope)
	}
	if respData.Port != nil {
		_ = d.Set("port", respData.Port)
	}
	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}
	if respData.ParamTemplateId != nil {
		_ = d.Set("param_template_id", respData.ParamTemplateId)
	}
	if respData.BelongMember != nil && respData.BelongMember.MemberId != nil {
		_ = d.Set("belong_member_id", respData.BelongMember.MemberId)
	}

	return nil
}

func resourceTencentCloudFwmNatAclRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_nat_acl_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = fwmv20250611.NewModifyNatAclRuleRequest()
	)

	groupId, ruleId, direction, err := parseFwmNatAclRuleId(d.Id())
	if err != nil {
		return err
	}

	request.GroupId = helper.String(groupId)

	rule := fwmv20250611.NatAclRule{}
	rule.RuleId = helper.String(ruleId)
	rule.Direction = helper.Int64(direction)
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
	if v, ok := d.GetOk("protocol"); ok {
		rule.Protocol = helper.String(v.(string))
	}
	if v, ok := d.GetOk("rule_action"); ok {
		rule.RuleAction = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("sequence"); ok {
		rule.OrderIndex = helper.Int64(int64(v.(int)))
	}
	if v, ok := d.GetOk("scope"); ok {
		rule.Scope = helper.String(v.(string))
	}
	if v, ok := d.GetOk("port"); ok {
		rule.Port = helper.String(v.(string))
	}
	if v, ok := d.GetOk("description"); ok {
		rule.Description = helper.String(v.(string))
	}
	if v, ok := d.GetOk("param_template_id"); ok {
		rule.ParamTemplateId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("belong_member_id"); ok {
		rule.BelongMemberId = helper.String(v.(string))
	}
	request.Rule = &rule

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().ModifyNatAclRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Modify fwm nat acl rule failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s modify fwm nat acl rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudFwmNatAclRuleRead(d, meta)
}

func resourceTencentCloudFwmNatAclRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_nat_acl_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = fwmv20250611.NewDeleteNatAclRuleRequest()
	)

	groupId, ruleId, _, err := parseFwmNatAclRuleId(d.Id())
	if err != nil {
		return err
	}

	request.GroupId = helper.String(groupId)
	request.RuleIds = []*string{helper.String(ruleId)}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().DeleteNatAclRuleWithContext(ctx, request)
		if e != nil {
			if sdkerr, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "ResourceNotFound" {
					return nil
				}
			}
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete fwm nat acl rule failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete fwm nat acl rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

// parseFwmNatAclRuleId splits the resource id <groupId>#<ruleId>#<direction>.
func parseFwmNatAclRuleId(id string) (groupId, ruleId string, direction int64, err error) {
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
