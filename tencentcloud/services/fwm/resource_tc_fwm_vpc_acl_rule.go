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

func ResourceTencentCloudFwmVpcAclRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudFwmVpcAclRuleCreate,
		Read:   resourceTencentCloudFwmVpcAclRuleRead,
		Update: resourceTencentCloudFwmVpcAclRuleUpdate,
		Delete: resourceTencentCloudFwmVpcAclRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Rule group ID, the parent resource reference, required by all four interfaces, immutable.",
			},

			"edge_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge firewall ID: ALL or CFWS-xxx.",
			},

			"fw_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Firewall instance ID: ALL, cfwg-xxx or ccn-xxx.",
			},

			"ip_version": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "IP version: ipv4 or ipv6. Not a Create input parameter, only returned by Describe.",
			},

			"source_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Source type: ip/url/template/instance/tag.",
			},

			"source_content": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Source address content.",
			},

			"dest_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Destination type: ip/url/template/instance/tag.",
			},

			"dest_content": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Destination address content.",
			},

			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Protocol: TCP/UDP/ICMP/ANY etc.",
			},

			"rule_action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Action: accept/drop/log (lowercase).",
			},

			"order_index": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Priority (starting from 1).",
			},

			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule description, required and up to 30 characters.",
			},

			"port": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Port, empty for ICMP protocol.",
			},

			"param_template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Port template ID.",
			},

			"belong_member_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule owner member account ID.",
			},

			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule ID.",
			},

			"sequence": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Priority (read-only, same source as order_index).",
			},

			"source_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source asset name.",
			},

			"dest_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Destination asset name.",
			},

			"param_template_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port template name.",
			},

			"fw_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Firewall group name.",
			},

			"invalid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether the rule is invalid: 0 valid, 1 invalid.",
			},

			"app_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Member AppId.",
			},

			"uin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Member Uin.",
			},

			"nickname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Member nickname.",
			},

			"member_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Member ID.",
			},

			"node_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Department.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule creation time.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule modification time.",
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

			"all_total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total count without filtering.",
			},
		},
	}
}

func resourceTencentCloudFwmVpcAclRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_vpc_acl_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = fwmv20250611.NewCreateVpcAclRuleRequest()
		response = fwmv20250611.NewCreateVpcAclRuleResponse()
		groupId  string
		ruleId   string
	)

	if v, ok := d.GetOk("group_id"); ok {
		request.GroupId = helper.String(v.(string))
		groupId = v.(string)
	}

	rule := fwmv20250611.VpcAclRule{}
	if v, ok := d.GetOk("source_type"); ok {
		rule.SourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_content"); ok {
		rule.SourceContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dest_type"); ok {
		rule.DestType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dest_content"); ok {
		rule.DestContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protocol"); ok {
		rule.Protocol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_action"); ok {
		rule.RuleAction = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_index"); ok {
		rule.OrderIndex = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("edge_id"); ok {
		rule.EdgeId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("fw_group_id"); ok {
		rule.FwGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		rule.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("port"); ok {
		rule.Port = helper.String(v.(string))
	}

	if v, ok := d.GetOk("param_template_id"); ok {
		rule.ParamTemplateId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("belong_member_id"); ok {
		rule.BelongMemberId = helper.String(v.(string))
	}

	request.Rules = append(request.Rules, &rule)

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().CreateVpcAclRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create fwm vpc_acl_rule failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create fwm vpc_acl_rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s create fwm vpc_acl_rule, d.Id()=%s", logId, d.Id())
	if response.Response.RuleIds == nil || len(response.Response.RuleIds) == 0 {
		return fmt.Errorf("create fwm vpc_acl_rule failed, RuleIds is empty")
	}

	ruleId = *response.Response.RuleIds[0]
	d.SetId(strings.Join([]string{groupId, ruleId}, tccommon.FILED_SP))
	return resourceTencentCloudFwmVpcAclRuleRead(d, meta)
}

func resourceTencentCloudFwmVpcAclRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_vpc_acl_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = NewFwmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	groupId := idSplit[0]
	ruleId := idSplit[1]

	respData, allTotalCount, err := service.DescribeFwmVpcAclRuleById(ctx, groupId, ruleId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[CRUD] fwm_vpc_acl_rule id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if respData.EdgeId != nil {
		_ = d.Set("edge_id", respData.EdgeId)
	}

	if respData.FwGroupId != nil {
		_ = d.Set("fw_group_id", respData.FwGroupId)
	}

	if respData.IpVersion != nil {
		_ = d.Set("ip_version", respData.IpVersion)
	}

	if respData.SourceType != nil {
		_ = d.Set("source_type", respData.SourceType)
	}

	if respData.SourceContent != nil {
		_ = d.Set("source_content", respData.SourceContent)
	}

	if respData.DestType != nil {
		_ = d.Set("dest_type", respData.DestType)
	}

	if respData.DestContent != nil {
		_ = d.Set("dest_content", respData.DestContent)
	}

	if respData.Protocol != nil {
		_ = d.Set("protocol", respData.Protocol)
	}

	if respData.RuleAction != nil {
		_ = d.Set("rule_action", respData.RuleAction)
	}

	if respData.Sequence != nil {
		_ = d.Set("order_index", respData.Sequence)
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.Port != nil {
		_ = d.Set("port", respData.Port)
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

	if respData.SourceName != nil {
		_ = d.Set("source_name", respData.SourceName)
	}

	if respData.DestName != nil {
		_ = d.Set("dest_name", respData.DestName)
	}

	if respData.ParamTemplateName != nil {
		_ = d.Set("param_template_name", respData.ParamTemplateName)
	}

	if respData.FwGroupName != nil {
		_ = d.Set("fw_group_name", respData.FwGroupName)
	}

	if respData.Invalid != nil {
		_ = d.Set("invalid", respData.Invalid)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.UpdateTime != nil {
		_ = d.Set("update_time", respData.UpdateTime)
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

	if respData.BelongMember != nil {
		if respData.BelongMember.AppId != nil {
			_ = d.Set("app_id", respData.BelongMember.AppId)
		}

		if respData.BelongMember.Uin != nil {
			_ = d.Set("uin", respData.BelongMember.Uin)
		}

		if respData.BelongMember.Nickname != nil {
			_ = d.Set("nickname", respData.BelongMember.Nickname)
		}

		if respData.BelongMember.MemberId != nil {
			_ = d.Set("member_id", respData.BelongMember.MemberId)
			_ = d.Set("belong_member_id", respData.BelongMember.MemberId)
		}

		if respData.BelongMember.NodeName != nil {
			_ = d.Set("node_name", respData.BelongMember.NodeName)
		}
	}

	if allTotalCount != nil {
		_ = d.Set("all_total_count", allTotalCount)
	}

	return nil
}

func resourceTencentCloudFwmVpcAclRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_vpc_acl_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	groupId := idSplit[0]
	ruleId := idSplit[1]

	needChange := false
	mutableArgs := []string{"protocol", "rule_action", "order_index", "description", "port", "param_template_id", "belong_member_id"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := fwmv20250611.NewModifyVpcAclRuleRequest()
		rule := fwmv20250611.VpcAclRule{}
		rule.RuleId = helper.String(ruleId)

		if v, ok := d.GetOk("source_type"); ok {
			rule.SourceType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("source_content"); ok {
			rule.SourceContent = helper.String(v.(string))
		}

		if v, ok := d.GetOk("dest_type"); ok {
			rule.DestType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("dest_content"); ok {
			rule.DestContent = helper.String(v.(string))
		}

		if v, ok := d.GetOk("protocol"); ok {
			rule.Protocol = helper.String(v.(string))
		}

		if v, ok := d.GetOk("rule_action"); ok {
			rule.RuleAction = helper.String(v.(string))
		}

		if v, ok := d.GetOk("order_index"); ok {
			rule.OrderIndex = helper.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOk("edge_id"); ok {
			rule.EdgeId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("fw_group_id"); ok {
			rule.FwGroupId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("description"); ok {
			rule.Description = helper.String(v.(string))
		}

		if v, ok := d.GetOk("port"); ok {
			rule.Port = helper.String(v.(string))
		}

		if v, ok := d.GetOk("param_template_id"); ok {
			rule.ParamTemplateId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("belong_member_id"); ok {
			rule.BelongMemberId = helper.String(v.(string))
		}

		request.Rule = &rule
		request.GroupId = helper.String(groupId)

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().ModifyVpcAclRuleWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify fwm vpc_acl_rule failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update fwm vpc_acl_rule failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudFwmVpcAclRuleRead(d, meta)
}

func resourceTencentCloudFwmVpcAclRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_vpc_acl_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = NewFwmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	groupId := idSplit[0]
	ruleId := idSplit[1]

	if err := service.DeleteFwmVpcAclRuleById(ctx, groupId, ruleId); err != nil {
		log.Printf("[CRITAL]%s delete fwm vpc_acl_rule failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
