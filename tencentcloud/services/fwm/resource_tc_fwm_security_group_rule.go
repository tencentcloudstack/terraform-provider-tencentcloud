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
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "ID of the enterprise security group rule group. Maps to `GroupId` of `CreateSecurityGroupRule`, `DescribeSecurityGroupRules`, `ModifySecurityGroupRule` and `DeleteSecurityGroupRule`. Exactly one of `group_id` and `group_name` must be specified.",
			},
			"group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Name of a new enterprise security group rule group to create for this rule. Maps to `GroupName` of `CreateRuleGroup` (the product is fixed to `enterprise_sg`). The group is created on apply and deleted on destroy. Exactly one of `group_id` and `group_name` must be specified.",
			},
			"ip_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "IP version of the rule. Maps to `IpVersion` of `SecurityGroupRule`.",
			},
			"source_content": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Source address of the rule. Maps to `SourceContent` of `SecurityGroupRule`. The API echoes it back as the int-encoded `SourceType`/`SourceId` pair, so it is not refreshed from the describe response.",
			},
			"source_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Source type of the rule. Maps to `SourceType` (string) of `SecurityGroupRule`. The describe response returns `SourceType` as an int (see the computed `target_type`), so this attribute is not refreshed from the response.",
			},
			"dest_content": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Destination address of the rule. Maps to `DestContent` of `SecurityGroupRule`. The describe response returns it as `TargetId`, so it is not refreshed from the response.",
			},
			"dest_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Destination type of the rule. Maps to `DestType` (string) of `SecurityGroupRule`. The describe response returns `TargetType` as an int (see the computed `target_type`), so this attribute is not refreshed from the response.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Protocol of the rule, for example `TCP`, `UDP`, `ICMP` or `ANY`. Maps to `Protocol` of `SecurityGroupRule`.",
			},
			"port": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Port or port range of the rule. Maps to `Port` of `SecurityGroupRule`.",
			},
			"service_template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the protocol port template used by the rule. Maps to `ServiceTemplateId` of `SecurityGroupRule`.",
			},
			"rule_action": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Action of the rule, for example `accept` or `drop`. Maps to `RuleAction` (string) of `SecurityGroupRule`. The describe response returns the action as an int `Strategy` (see the computed `strategy`), so this attribute is not refreshed from the response.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Description of the rule. Maps to `Description` of `SecurityGroupRule`. The describe response returns the remark in `Detail` (see the computed `detail`), so this attribute is not refreshed from the response.",
			},
			"order_index": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Priority of the rule, the smaller the value the higher the priority. Maps to `OrderIndex` of `SecurityGroupRule`.",
			},
			"scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Effective scope of the rule. Maps to `Scope` of `SecurityGroupRule`.",
			},
			"protocol_port_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Protocol port type of the rule. Maps to `ProtocolPortType` of `SecurityGroupRule`. The describe response uses a different encoding, so this attribute is not refreshed from the response.",
			},
			"belong_member_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the member account the rule belongs to. Maps to `BelongMemberId` of `SecurityGroupRule`. The describe response returns a `MemberInfo` object (see the computed `belong_member`), so this attribute is not refreshed from the response.",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the rule, returned by `CreateSecurityGroupRule` as `RuleIds[0]` and echoed by the describe response as `RuleId`.",
			},
			"source_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source rule content returned by the describe response as `SourceId`.",
			},
			"target_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Destination rule content returned by the describe response as `TargetId`.",
			},
			"target_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Destination rule type returned by the describe response as `TargetType`.",
			},
			"strategy": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Rule action returned by the describe response as `Strategy`.",
			},
			"detail": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule remark returned by the describe response as `Detail`.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region of the rule returned by the describe response as `Region`.",
			},
			"sou_instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source asset name returned by the describe response as `SouInstanceName`.",
			},
			"sou_public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source asset public IP returned by the describe response as `SouPublicIp`.",
			},
			"sou_private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source asset private IP returned by the describe response as `SouPrivateIp`.",
			},
			"sou_cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source asset CIDR returned by the describe response as `SouCidr`.",
			},
			"sou_parameter_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source template name returned by the describe response as `SouParameterName`.",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Destination asset name returned by the describe response as `InstanceName`.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Destination asset public IP returned by the describe response as `PublicIp`.",
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Destination asset private IP returned by the describe response as `PrivateIp`.",
			},
			"cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Destination asset CIDR returned by the describe response as `Cidr`.",
			},
			"parameter_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Destination template name returned by the describe response as `ParameterName`.",
			},
			"protocol_port_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Protocol port template name returned by the describe response as `ProtocolPortName`.",
			},
			"is_new": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether the rule was changed recently, `0` means no and `1` means yes. Returned by the describe response as `IsNew`.",
			},
			"rule_uuid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Numeric rule ID (uuid), equal to `RuleId`. Returned by the describe response as `Id`.",
			},
			"dns_parse_count": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Description: "Statistics of the IPs resolved from the domain name of the rule. Returned by the describe response as `DnsParseCount`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"valid_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of IPs successfully delivered. Maps to `ValidCount` of `SgDnsParseCount`.",
						},
						"invalid_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of IPs not delivered. Maps to `InvalidCount` of `SgDnsParseCount`.",
						},
					},
				},
			},
			"belong_member": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Description: "Member account the rule belongs to. Returned by the describe response as `BelongMember`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AppId of the member. Maps to `AppId` of `MemberInfo`.",
						},
						"uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Uin of the member. Maps to `Uin` of `MemberInfo`.",
						},
						"nickname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Nickname of the member. Maps to `Nickname` of `MemberInfo`.",
						},
						"member_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the member, empty in personal scenarios. Maps to `MemberId` of `MemberInfo`.",
						},
						"node_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Department the member belongs to. Maps to `NodeName` of `MemberInfo`.",
						},
					},
				},
			},
		},
	}
}

func buildFwmSecurityGroupRuleFromSchema(d *schema.ResourceData) *fwmv20250611.SecurityGroupRule {
	rule := fwmv20250611.SecurityGroupRule{}

	if v, ok := d.GetOk("ip_version"); ok {
		rule.IpVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_content"); ok {
		rule.SourceContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_type"); ok {
		rule.SourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dest_content"); ok {
		rule.DestContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dest_type"); ok {
		rule.DestType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protocol"); ok {
		rule.Protocol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("port"); ok {
		rule.Port = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_template_id"); ok {
		rule.ServiceTemplateId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_action"); ok {
		rule.RuleAction = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		rule.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_index"); ok {
		rule.OrderIndex = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("scope"); ok {
		rule.Scope = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protocol_port_type"); ok {
		rule.ProtocolPortType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("belong_member_id"); ok {
		rule.BelongMemberId = helper.String(v.(string))
	}

	return &rule
}

func parseFwmSecurityGroupRuleId(id string) (groupId, ruleId string, err error) {
	idSplit := strings.Split(id, tccommon.FILED_SP)
	if len(idSplit) != 2 || idSplit[0] == "" || idSplit[1] == "" {
		return "", "", fmt.Errorf("id is broken,%s", id)
	}

	return idSplit[0], idSplit[1], nil
}

func resourceTencentCloudFwmSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_security_group_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		client  = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		groupId string
		ruleId  string
	)

	_, hasGroupId := d.GetOk("group_id")
	_, hasGroupName := d.GetOk("group_name")
	if !hasGroupId && !hasGroupName {
		return fmt.Errorf("tencentcloud_fwm_security_group_rule create failed: either group_id or group_name must be specified")
	}

	if hasGroupId && hasGroupName {
		return fmt.Errorf("tencentcloud_fwm_security_group_rule create failed: group_id and group_name are mutually exclusive")
	}

	if hasGroupName {
		groupName := d.Get("group_name").(string)
		groupRequest := fwmv20250611.NewCreateRuleGroupRequest()
		groupRequest.GroupName = helper.String(groupName)
		groupRequest.Product = helper.String("enterprise_sg")

		var groupResponse *fwmv20250611.CreateRuleGroupResponse
		groupErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := client.UseFwmV20250611Client().CreateRuleGroupWithContext(ctx, groupRequest)
			if e != nil {
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, groupRequest.GetAction(), groupRequest.ToJsonString(), result.ToJsonString())

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("CreateRuleGroup of tencentcloud_fwm_security_group_rule returned an empty response."))
			}

			groupResponse = result
			return nil
		})

		if groupErr != nil {
			log.Printf("[CRITAL]%s create rule group of tencentcloud_fwm_security_group_rule failed, reason:%+v", logId, groupErr)
			return groupErr
		}

		if groupResponse.Response.GroupId == nil || *groupResponse.Response.GroupId == "" {
			log.Printf("[CRITAL]%s create rule group of tencentcloud_fwm_security_group_rule failed, logId:%s, GroupId is empty", logId, logId)
			return fmt.Errorf("CreateRuleGroup of tencentcloud_fwm_security_group_rule returned an empty GroupId")
		}

		groupId = *groupResponse.Response.GroupId
	} else {
		groupId = d.Get("group_id").(string)
	}

	request := fwmv20250611.NewCreateSecurityGroupRuleRequest()
	request.GroupId = helper.String(groupId)
	request.Rules = []*fwmv20250611.SecurityGroupRule{buildFwmSecurityGroupRuleFromSchema(d)}

	var response *fwmv20250611.CreateSecurityGroupRuleResponse
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := client.UseFwmV20250611Client().CreateSecurityGroupRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("CreateSecurityGroupRule of tencentcloud_fwm_security_group_rule returned an empty response."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create tencentcloud_fwm_security_group_rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if len(response.Response.RuleIds) == 0 || response.Response.RuleIds[0] == nil || *response.Response.RuleIds[0] == "" {
		log.Printf("[CRITAL]%s create tencentcloud_fwm_security_group_rule failed, logId:%s, id:%s, RuleIds is empty", logId, logId, d.Id())
		return fmt.Errorf("CreateSecurityGroupRule of tencentcloud_fwm_security_group_rule returned an empty RuleIds")
	}

	ruleId = *response.Response.RuleIds[0]
	log.Printf("[DEBUG]%s create tencentcloud_fwm_security_group_rule success, logId:%s, id:%s, groupId:%s, ruleId:%s", logId, logId, d.Id(), groupId, ruleId)

	d.SetId(strings.Join([]string{groupId, ruleId}, tccommon.FILED_SP))
	_ = d.Set("rule_id", ruleId)
	_ = d.Set("group_id", groupId)

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

	groupId, ruleId, err := parseFwmSecurityGroupRuleId(d.Id())
	if err != nil {
		return err
	}

	respData, err := service.DescribeFwmSecurityGroupRuleById(ctx, groupId, ruleId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[CRUD] tencentcloud_fwm_security_group_rule id=%s", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("group_id", groupId)
	_ = d.Set("rule_id", ruleId)

	if respData.IpVersion != nil {
		_ = d.Set("ip_version", respData.IpVersion)
	}

	if respData.Protocol != nil {
		_ = d.Set("protocol", respData.Protocol)
	}

	if respData.Port != nil {
		_ = d.Set("port", respData.Port)
	}

	if respData.ServiceTemplateId != nil {
		_ = d.Set("service_template_id", respData.ServiceTemplateId)
	}

	if respData.Scope != nil {
		_ = d.Set("scope", respData.Scope)
	}

	if respData.OrderIndex != nil {
		_ = d.Set("order_index", respData.OrderIndex)
	}

	if respData.SourceId != nil {
		_ = d.Set("source_id", respData.SourceId)
	}

	if respData.TargetId != nil {
		_ = d.Set("target_id", respData.TargetId)
	}

	if respData.TargetType != nil {
		_ = d.Set("target_type", respData.TargetType)
	}

	if respData.Strategy != nil {
		_ = d.Set("strategy", respData.Strategy)
	}

	if respData.Detail != nil {
		_ = d.Set("detail", respData.Detail)
	}

	if respData.Region != nil {
		_ = d.Set("region", respData.Region)
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

	if respData.ProtocolPortName != nil {
		_ = d.Set("protocol_port_name", respData.ProtocolPortName)
	}

	if respData.IsNew != nil {
		_ = d.Set("is_new", respData.IsNew)
	}

	if respData.Id != nil {
		_ = d.Set("rule_uuid", respData.Id)
	}

	if respData.DnsParseCount != nil {
		dnsParseCount := map[string]interface{}{}
		if respData.DnsParseCount.ValidCount != nil {
			dnsParseCount["valid_count"] = respData.DnsParseCount.ValidCount
		}

		if respData.DnsParseCount.InvalidCount != nil {
			dnsParseCount["invalid_count"] = respData.DnsParseCount.InvalidCount
		}

		_ = d.Set("dns_parse_count", []map[string]interface{}{dnsParseCount})
	}

	if respData.BelongMember != nil {
		belongMember := map[string]interface{}{}
		if respData.BelongMember.AppId != nil {
			belongMember["app_id"] = respData.BelongMember.AppId
		}

		if respData.BelongMember.Uin != nil {
			belongMember["uin"] = respData.BelongMember.Uin
		}

		if respData.BelongMember.Nickname != nil {
			belongMember["nickname"] = respData.BelongMember.Nickname
		}

		if respData.BelongMember.MemberId != nil {
			belongMember["member_id"] = respData.BelongMember.MemberId
		}

		if respData.BelongMember.NodeName != nil {
			belongMember["node_name"] = respData.BelongMember.NodeName
		}

		_ = d.Set("belong_member", []map[string]interface{}{belongMember})
	}

	return nil
}

func resourceTencentCloudFwmSecurityGroupRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_security_group_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId  = tccommon.GetLogId(tccommon.ContextNil)
		ctx    = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		client = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	)

	groupId, ruleId, err := parseFwmSecurityGroupRuleId(d.Id())
	if err != nil {
		return err
	}

	needChange := false
	mutableArgs := []string{
		"ip_version", "source_content", "source_type", "dest_content", "dest_type",
		"protocol", "port", "service_template_id", "rule_action", "description",
		"order_index", "scope", "protocol_port_type", "belong_member_id",
	}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := fwmv20250611.NewModifySecurityGroupRuleRequest()
		request.GroupId = helper.String(groupId)
		rule := buildFwmSecurityGroupRuleFromSchema(d)
		rule.RuleId = helper.String(ruleId)
		request.Rule = rule

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := client.UseFwmV20250611Client().ModifySecurityGroupRuleWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("ModifySecurityGroupRule of tencentcloud_fwm_security_group_rule returned an empty response."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update tencentcloud_fwm_security_group_rule failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudFwmSecurityGroupRuleRead(d, meta)
}

func resourceTencentCloudFwmSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_security_group_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId  = tccommon.GetLogId(tccommon.ContextNil)
		ctx    = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		client = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	)

	groupId, ruleId, err := parseFwmSecurityGroupRuleId(d.Id())
	if err != nil {
		return err
	}

	request := fwmv20250611.NewDeleteSecurityGroupRuleRequest()
	request.GroupId = helper.String(groupId)
	request.RuleIds = []*string{helper.String(ruleId)}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := client.UseFwmV20250611Client().DeleteSecurityGroupRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("DeleteSecurityGroupRule of tencentcloud_fwm_security_group_rule returned an empty response."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete tencentcloud_fwm_security_group_rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if _, ok := d.GetOk("group_name"); ok {
		groupRequest := fwmv20250611.NewDeleteRuleGroupRequest()
		groupRequest.GroupIds = []*string{helper.String(groupId)}

		groupErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := client.UseFwmV20250611Client().DeleteRuleGroupWithContext(ctx, groupRequest)
			if e != nil {
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, groupRequest.GetAction(), groupRequest.ToJsonString(), result.ToJsonString())

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("DeleteRuleGroup of tencentcloud_fwm_security_group_rule returned an empty response."))
			}

			return nil
		})

		if groupErr != nil {
			log.Printf("[WARN]%s delete rule group[%s] of tencentcloud_fwm_security_group_rule failed, reason:%+v", logId, groupId, groupErr)
		}
	}

	return nil
}
