package fwm

import (
	"context"
	"fmt"
	"log"
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
				Description: "Rule direction: 0-outbound, 1-inbound.",
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

			"order_index": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Priority, >=1, per-direction numbering.",
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

			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule ID, part of composite id.",
			},

			"sequence": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Sequence, display mirror of order_index.",
			},

			"scope_desc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Effective range description (language-dependent display).",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Update time.",
			},

			"source_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source asset name (display only).",
			},

			"target_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Target asset name (display only).",
			},

			"param_template_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port template name (display only).",
			},

			"fw_ins_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Firewall instance ID.",
			},

			"country_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Country name.",
			},

			"city_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "City name.",
			},

			"country_code": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Country code.",
			},

			"city_code": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "City code.",
			},

			"country_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Country key.",
			},

			"city_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "City key.",
			},

			"is_region": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether region rule: 0-no, 1-yes.",
			},

			"cloud_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud vendor code.",
			},

			"is_cloud": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether cloud vendor rule: 0-no, 1-yes.",
			},

			"invalid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether rule is invalid: 0-valid, 1-invalid.",
			},

			"app_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Member AppId (from BelongMember).",
			},

			"uin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Member Uin (from BelongMember).",
			},

			"nickname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Member nickname (from BelongMember).",
			},

			"member_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Member ID (from BelongMember).",
			},

			"node_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Member department (from BelongMember).",
			},

			"dns_parse_cnt": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Domain count.",
			},

			"hit_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Hit count.",
			},

			"last_hit_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last hit time.",
			},
		},
	}
}

func resourceTencentCloudFwmNatAclRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_nat_acl_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = fwmv20250611.NewCreateNatAclRuleRequest()
		response = fwmv20250611.NewCreateNatAclRuleResponse()
		groupId  string
		ruleId   string
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
	if v, ok := d.GetOkExists("order_index"); ok {
		rule.OrderIndex = helper.Int64(int64(v.(int)))
	}
	if v, ok := d.GetOk("scope"); ok {
		rule.Scope = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("direction"); ok {
		rule.Direction = helper.Int64(int64(v.(int)))
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
	d.SetId(strings.Join([]string{groupId, ruleId}, tccommon.FILED_SP))
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

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	groupId := idSplit[0]
	ruleId := idSplit[1]

	direction := int64(d.Get("direction").(int))

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
		log.Printf("[CRUD] fwm_nat_acl_rule id=%s", d.Id())
		d.SetId("")
		return nil
	}

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
		_ = d.Set("order_index", respData.Sequence)
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
	if respData.BelongMember != nil {
		if respData.BelongMember.MemberId != nil {
			_ = d.Set("belong_member_id", respData.BelongMember.MemberId)
			_ = d.Set("member_id", respData.BelongMember.MemberId)
		}
		if respData.BelongMember.AppId != nil {
			_ = d.Set("app_id", respData.BelongMember.AppId)
		}
		if respData.BelongMember.Uin != nil {
			_ = d.Set("uin", respData.BelongMember.Uin)
		}
		if respData.BelongMember.Nickname != nil {
			_ = d.Set("nickname", respData.BelongMember.Nickname)
		}
		if respData.BelongMember.NodeName != nil {
			_ = d.Set("node_name", respData.BelongMember.NodeName)
		}
	}
	if respData.ScopeDesc != nil {
		_ = d.Set("scope_desc", respData.ScopeDesc)
	}
	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}
	if respData.UpdateTime != nil {
		_ = d.Set("update_time", respData.UpdateTime)
	}
	if respData.SourceName != nil {
		_ = d.Set("source_name", respData.SourceName)
	}
	if respData.TargetName != nil {
		_ = d.Set("target_name", respData.TargetName)
	}
	if respData.ParamTemplateName != nil {
		_ = d.Set("param_template_name", respData.ParamTemplateName)
	}
	if respData.FwInsId != nil {
		_ = d.Set("fw_ins_id", respData.FwInsId)
	}
	if respData.CountryName != nil {
		_ = d.Set("country_name", respData.CountryName)
	}
	if respData.CityName != nil {
		_ = d.Set("city_name", respData.CityName)
	}
	if respData.CountryCode != nil {
		_ = d.Set("country_code", respData.CountryCode)
	}
	if respData.CityCode != nil {
		_ = d.Set("city_code", respData.CityCode)
	}
	if respData.CountryKey != nil {
		_ = d.Set("country_key", respData.CountryKey)
	}
	if respData.CityKey != nil {
		_ = d.Set("city_key", respData.CityKey)
	}
	if respData.IsRegion != nil {
		_ = d.Set("is_region", respData.IsRegion)
	}
	if respData.CloudCode != nil {
		_ = d.Set("cloud_code", respData.CloudCode)
	}
	if respData.IsCloud != nil {
		_ = d.Set("is_cloud", respData.IsCloud)
	}
	if respData.Invalid != nil {
		_ = d.Set("invalid", respData.Invalid)
	}
	if respData.DnsParseCnt != nil {
		_ = d.Set("dns_parse_cnt", respData.DnsParseCnt)
	}
	if respData.HitCount != nil {
		_ = d.Set("hit_count", respData.HitCount)
	}
	if respData.LastHitTime != nil {
		_ = d.Set("last_hit_time", respData.LastHitTime)
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

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	groupId := idSplit[0]
	ruleId := idSplit[1]

	request.GroupId = helper.String(groupId)

	rule := fwmv20250611.NatAclRule{}
	rule.RuleId = helper.String(ruleId)
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
	if v, ok := d.GetOkExists("order_index"); ok {
		rule.OrderIndex = helper.Int64(int64(v.(int)))
	}
	if v, ok := d.GetOk("scope"); ok {
		rule.Scope = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("direction"); ok {
		rule.Direction = helper.Int64(int64(v.(int)))
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

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	groupId := idSplit[0]
	ruleId := idSplit[1]

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
