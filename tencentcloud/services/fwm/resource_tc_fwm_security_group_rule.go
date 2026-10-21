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
				Description: "Rule group ID, the parent rule group (Product=enterprise_sg) to which the rule belongs.",
			},

			"ip_version": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "IP version, optional values: ipv4, ipv6.",
			},

			"order_index": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Rule priority order index, global single sequence within the group, >=1.",
			},

			"source_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Source type, optional values: net, vpc, subnet, cvm, clb, eni, cdb, template, tag, region, resourcegroup, addressbook, instance.",
			},

			"source_content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Source content, such as IP/CIDR or instance id.",
			},

			"dest_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Destination type, optional values: net, vpc, subnet, cvm, clb, eni, cdb, template, tag, region, resourcegroup, addressbook, instance.",
			},

			"dest_content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Destination content, such as IP/CIDR or instance id.",
			},

			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Protocol, optional values: TCP, ANY, ICMP, UDP. ANY means all.",
			},

			"port": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Port, such as 80 or ALL.",
			},

			"rule_action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule action, optional values: accept, drop, log.",
			},

			"scope": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Rule effective scope, optional values: SG, LH.",
			},

			"service_template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Service template id.",
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
				Description: "Member id to which the rule belongs.",
			},

			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule id.",
			},

			"source_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source id (read-only display).",
			},

			"target_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Target id (read-only display).",
			},

			"strategy": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Rule strategy int value (read-only display).",
			},

			"detail": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule detail (read-only display).",
			},

			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region (read-only display).",
			},

			"protocol_port_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Protocol port name (read-only display).",
			},

			"sou_instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source instance name (read-only display).",
			},

			"sou_public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source public ip (read-only display).",
			},

			"sou_private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source private ip (read-only display).",
			},

			"sou_cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source cidr (read-only display).",
			},

			"sou_parameter_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source parameter name (read-only display).",
			},

			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Destination instance name (read-only display).",
			},

			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Destination public ip (read-only display).",
			},

			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Destination private ip (read-only display).",
			},

			"cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Destination cidr (read-only display).",
			},

			"parameter_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Destination parameter name (read-only display).",
			},

			"dns_parse_count": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"valid_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Valid ip count.",
						},
						"invalid_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Invalid ip count.",
						},
					},
				},
				Description: "Dns parse count (read-only display).",
			},

			"is_new": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether the rule has been changed recently (read-only display).",
			},

			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Internal primary key id (read-only display).",
			},

			"all_total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "All total count (read-only display).",
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
	)

	if v, ok := d.GetOk("group_id"); ok {
		request.GroupId = helper.String(v.(string))
	}

	rule.IpVersion = helper.String(d.Get("ip_version").(string))
	rule.OrderIndex = helper.Int64(int64(d.Get("order_index").(int)))
	rule.SourceType = helper.String(d.Get("source_type").(string))
	rule.SourceContent = helper.String(d.Get("source_content").(string))
	rule.DestType = helper.String(d.Get("dest_type").(string))
	rule.DestContent = helper.String(d.Get("dest_content").(string))
	rule.Protocol = helper.String(d.Get("protocol").(string))
	rule.Port = helper.String(d.Get("port").(string))
	rule.RuleAction = helper.String(d.Get("rule_action").(string))
	rule.Scope = helper.String(d.Get("scope").(string))

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

	request.Rules = []*fwmv20250611.SecurityGroupRule{&rule}

	var ruleId string
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
		log.Printf("[CRUD] fwm_security_group_rule id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if respData.RuleId != nil {
		_ = d.Set("rule_id", respData.RuleId)
	}

	if respData.IpVersion != nil {
		_ = d.Set("ip_version", respData.IpVersion)
	}

	if respData.OrderIndex != nil {
		_ = d.Set("order_index", respData.OrderIndex)
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
		_ = d.Set("strategy", respData.Strategy)
	}

	if respData.Detail != nil {
		_ = d.Set("description", respData.Detail)
		_ = d.Set("detail", respData.Detail)
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

	if respData.Region != nil {
		_ = d.Set("region", respData.Region)
	}

	if respData.ProtocolPortName != nil {
		_ = d.Set("protocol_port_name", respData.ProtocolPortName)
	}

	if respData.SouInstanceName != nil {
		_ = d.Set("sou_instance_name", respData.SouInstanceName)
	}

	if respData.SouPublicIp != nil {
		_ = d.Set("sou_public_ip", respData.SouPublicIp)
	}

	if respData.SouPrivateIp != nil {
		_ = d.Set("sou_private_ip", respData.SouPrivateIp)
	}

	if respData.SouCidr != nil {
		_ = d.Set("sou_cidr", respData.SouCidr)
	}

	if respData.SouParameterName != nil {
		_ = d.Set("sou_parameter_name", respData.SouParameterName)
	}

	if respData.InstanceName != nil {
		_ = d.Set("instance_name", respData.InstanceName)
	}

	if respData.PublicIp != nil {
		_ = d.Set("public_ip", respData.PublicIp)
	}

	if respData.PrivateIp != nil {
		_ = d.Set("private_ip", respData.PrivateIp)
	}

	if respData.Cidr != nil {
		_ = d.Set("cidr", respData.Cidr)
	}

	if respData.ParameterName != nil {
		_ = d.Set("parameter_name", respData.ParameterName)
	}

	if respData.IsNew != nil {
		_ = d.Set("is_new", respData.IsNew)
	}

	if respData.Id != nil {
		_ = d.Set("id", respData.Id)
	}

	if respData.DnsParseCount != nil {
		dnsParseCountMap := map[string]interface{}{}
		if respData.DnsParseCount.ValidCount != nil {
			dnsParseCountMap["valid_count"] = respData.DnsParseCount.ValidCount
		}

		if respData.DnsParseCount.InvalidCount != nil {
			dnsParseCountMap["invalid_count"] = respData.DnsParseCount.InvalidCount
		}
		_ = d.Set("dns_parse_count", []interface{}{dnsParseCountMap})
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

	mutableArgs := []string{"order_index", "source_content", "dest_content", "protocol", "port", "service_template_id", "rule_action", "description", "protocol_port_type", "belong_member_id"}
	needChange := false
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := fwmv20250611.NewModifySecurityGroupRuleRequest()
		request.GroupId = helper.String(groupId)

		rule := fwmv20250611.SecurityGroupRule{}
		rule.RuleId = helper.String(ruleId)
		rule.IpVersion = helper.String(d.Get("ip_version").(string))
		rule.OrderIndex = helper.Int64(int64(d.Get("order_index").(int)))
		rule.SourceType = helper.String(d.Get("source_type").(string))
		rule.SourceContent = helper.String(d.Get("source_content").(string))
		rule.DestType = helper.String(d.Get("dest_type").(string))
		rule.DestContent = helper.String(d.Get("dest_content").(string))
		rule.Protocol = helper.String(d.Get("protocol").(string))
		rule.Port = helper.String(d.Get("port").(string))
		rule.RuleAction = helper.String(d.Get("rule_action").(string))
		rule.Scope = helper.String(d.Get("scope").(string))

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

		request.Rule = &rule

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
