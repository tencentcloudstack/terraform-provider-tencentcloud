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

var sourceTypeInt2Str = map[int64]string{
	0:   "net",
	1:   "vpc",
	2:   "subnet",
	3:   "cvm",
	4:   "clb",
	5:   "eni",
	6:   "cdb",
	7:   "template",
	8:   "tag",
	9:   "region",
	100: "resourcegroup",
}

var ruleActionInt2Str = map[int64]string{
	0: "log",
	1: "drop",
	2: "accept",
}

func reverseSourceType(i int64) string {
	if v, ok := sourceTypeInt2Str[i]; ok {
		return v
	}
	log.Printf("[WARN] fwm_security_group_rule unknown SourceType int %d, keep decimal string", i)
	return strconv.FormatInt(i, 10)
}

func reverseRuleAction(i int64) string {
	if v, ok := ruleActionInt2Str[i]; ok {
		return v
	}
	log.Printf("[WARN] fwm_security_group_rule unknown Strategy int %d, keep decimal string", i)
	return strconv.FormatInt(i, 10)
}

func ResourceTencentCloudFwmSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudFwmSecurityGroupRuleCreate,
		Read:   resourceTencentCloudFwmSecurityGroupRuleRead,
		Update: resourceTencentCloudFwmSecurityGroupRuleUpdate,
		Delete: resourceTencentCloudFwmSecurityGroupRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Rule group ID, the parent rule group (Product=enterprise_sg) to which the rule belongs.",
			},

			"ip_version": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "IP version: ipv4, ipv6.",
			},

			"source_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Source type: net, vpc, subnet, cvm, clb, eni, cdb, template, tag, region, resourcegroup, addressbook, instance.",
			},

			"source_content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Source content, such as IP, CIDR or instance ID.",
			},

			"dest_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Destination type: net, vpc, subnet, cvm, clb, eni, cdb, template, tag, region, resourcegroup, addressbook, instance.",
			},

			"dest_content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Destination content, such as IP, CIDR or instance ID.",
			},

			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Protocol: TCP, UDP, ICMP, ANY. ANY means all protocols.",
			},

			"port": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Port, such as 80, 1-65535 or ALL.",
			},

			"rule_action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule action: accept, drop, log.",
			},

			"scope": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Rule effective scope: SG, LH.",
			},

			"service_template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Service template ID.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule description.",
			},

			"protocol_port_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Protocol port type.",
			},

			"belong_member_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the member account that the rule belongs to.",
			},

			// computed
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule ID.",
			},

			"source_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source ID.",
			},

			"target_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Destination ID.",
			},

			"sequence": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Priority (read-only, same source as order_index).",
			},
		},
	}
}

func resourceTencentCloudFwmSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_security_group_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = fwmv20250611.NewCreateSecurityGroupRuleRequest()
		rule    = fwmv20250611.SecurityGroupRule{}
		ruleId  string
	)

	if v, ok := d.GetOk("group_id"); ok {
		request.GroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ip_version"); ok {
		rule.IpVersion = helper.String(v.(string))
	}

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

	if v, ok := d.GetOk("port"); ok {
		rule.Port = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_action"); ok {
		rule.RuleAction = helper.String(v.(string))
	}

	if v, ok := d.GetOk("scope"); ok {
		rule.Scope = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_template_id"); ok {
		rule.ServiceTemplateId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		rule.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("protocol_port_type"); ok {
		rule.ProtocolPortType = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("belong_member_id"); ok {
		rule.BelongMemberId = helper.String(v.(string))
	}

	rule.OrderIndex = helper.Int64(int64(1))
	request.Rules = []*fwmv20250611.SecurityGroupRule{&rule}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().CreateSecurityGroupRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || len(result.Response.RuleIds) == 0 || result.Response.RuleIds[0] == nil || *result.Response.RuleIds[0] == "" {
			return resource.NonRetryableError(fmt.Errorf("Create fwm_security_group_rule failed, RuleIds is empty, logId: %s, d.Id(): %s", logId, d.Id()))
		}

		ruleId = *result.Response.RuleIds[0]
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create fwm_security_group_rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s fwm_security_group_rule created, logId: %s, d.Id(): %s", logId, logId, d.Id())

	groupId := d.Get("group_id").(string)
	d.SetId(strings.Join([]string{groupId, ruleId}, tccommon.FILED_SP))
	return resourceTencentCloudFwmSecurityGroupRuleRead(d, meta)
}

func resourceTencentCloudFwmSecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_security_group_rule.read")()
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

	respData, err := service.DescribeFwmSecurityGroupRuleById(ctx, groupId, ruleId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_fwm_security_group_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("group_id", groupId)

	if respData.RuleId != nil {
		_ = d.Set("rule_id", respData.RuleId)
	}

	if respData.IpVersion != nil {
		_ = d.Set("ip_version", respData.IpVersion)
	}

	if respData.OrderIndex != nil {
		_ = d.Set("sequence", respData.OrderIndex)
	}

	if respData.SourceType != nil {
		_ = d.Set("source_type", reverseSourceType(*respData.SourceType))
	}

	if respData.SourceId != nil {
		_ = d.Set("source_content", respData.SourceId)
		_ = d.Set("source_id", respData.SourceId)
	}

	if respData.TargetType != nil {
		_ = d.Set("dest_type", reverseSourceType(*respData.TargetType))
	}

	if respData.TargetId != nil {
		_ = d.Set("dest_content", respData.TargetId)
		_ = d.Set("target_id", respData.TargetId)
	}

	if respData.Protocol != nil {
		_ = d.Set("protocol", respData.Protocol)
	}

	if respData.Port != nil {
		_ = d.Set("port", respData.Port)
	}

	if respData.Strategy != nil {
		_ = d.Set("rule_action", reverseRuleAction(*respData.Strategy))
	}

	if respData.Detail != nil {
		_ = d.Set("description", respData.Detail)
	}

	if respData.Scope != nil {
		_ = d.Set("scope", respData.Scope)
	}

	if respData.ServiceTemplateId != nil {
		_ = d.Set("service_template_id", respData.ServiceTemplateId)
	}

	if respData.BelongMember != nil && respData.BelongMember.MemberId != nil && *respData.BelongMember.MemberId != "" {
		_ = d.Set("belong_member_id", respData.BelongMember.MemberId)
	}

	return nil
}

func resourceTencentCloudFwmSecurityGroupRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_security_group_rule.update")()
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

	mutableArgs := []string{"source_content", "dest_content", "protocol", "port", "service_template_id", "rule_action", "description", "protocol_port_type", "belong_member_id"}
	needChange := false
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := fwmv20250611.NewModifySecurityGroupRuleRequest()
		rule := fwmv20250611.SecurityGroupRule{}

		if v, ok := d.GetOk("ip_version"); ok {
			rule.IpVersion = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("sequence"); ok {
			rule.OrderIndex = helper.Int64(int64(v.(int)))
		} else {
			rule.OrderIndex = helper.Int64(1)
		}

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

		if v, ok := d.GetOk("port"); ok {
			rule.Port = helper.String(v.(string))
		}

		if v, ok := d.GetOk("rule_action"); ok {
			rule.RuleAction = helper.String(v.(string))
		}

		if v, ok := d.GetOk("scope"); ok {
			rule.Scope = helper.String(v.(string))
		}

		if v, ok := d.GetOk("service_template_id"); ok {
			rule.ServiceTemplateId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("description"); ok {
			rule.Description = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("protocol_port_type"); ok {
			rule.ProtocolPortType = helper.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOk("belong_member_id"); ok {
			rule.BelongMemberId = helper.String(v.(string))
		}

		rule.RuleId = helper.String(ruleId)
		request.Rule = &rule
		request.GroupId = helper.String(groupId)
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().ModifySecurityGroupRuleWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify fwm_security_group_rule failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update fwm_security_group_rule failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudFwmSecurityGroupRuleRead(d, meta)
}

func resourceTencentCloudFwmSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_security_group_rule.delete")()
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

	if err := service.DeleteFwmSecurityGroupRuleById(ctx, groupId, ruleId); err != nil {
		log.Printf("[CRITAL]%s delete fwm_security_group_rule failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
