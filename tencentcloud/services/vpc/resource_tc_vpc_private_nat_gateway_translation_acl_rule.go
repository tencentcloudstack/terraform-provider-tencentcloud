package vpc

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpcv20170312 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcPrivateNatGatewayTranslationAclRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcPrivateNatGatewayTranslationAclRuleCreate,
		Read:   resourceTencentCloudVpcPrivateNatGatewayTranslationAclRuleRead,
		Update: resourceTencentCloudVpcPrivateNatGatewayTranslationAclRuleUpdate,
		Delete: resourceTencentCloudVpcPrivateNatGatewayTranslationAclRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"nat_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique ID of the private NAT gateway, in the format: `intranat-xxxxxxxx`.",
			},

			"translation_direction": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The target of the translation rule, optional value: LOCAL.",
			},

			"translation_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of translation rule, optional values: NETWORK_LAYER, TRANSPORT_LAYER. Corresponding to layer 3 and layer 4 respectively.",
			},

			"translation_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The mapped IP address. When the translation rule type is layer 4, it represents an IP pool.",
			},

			"translation_acl_rules": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Access control list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ACL protocol type, optional values: `ALL`, `TCP`, `UDP`.",
						},
						"source_port": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Source port.",
						},
						"source_cidr": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Source address. Supports `ip` or `cidr` format `xxx.xxx.xxx.000/xx`.",
						},
						"destination_port": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Destination port.",
						},
						"destination_cidr": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Destination address.",
						},
						"acl_rule_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ACL rule ID.",
						},
						"action": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Whether to match.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ACL rule description.",
						},
					},
				},
			},

			"original_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The original IP address before mapping. Valid when the translation rule type is layer 3.",
			},
		},
	}
}

func resourceTencentCloudVpcPrivateNatGatewayTranslationAclRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_private_nat_gateway_translation_acl_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                = tccommon.GetLogId(tccommon.ContextNil)
		ctx                  = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request              = vpcv20170312.NewCreatePrivateNatGatewayTranslationAclRuleRequest()
		response             = vpcv20170312.NewCreatePrivateNatGatewayTranslationAclRuleResponse()
		natGatewayId         string
		translationDirection string
		translationType      string
		translationIp        string
		originalIp           string
		aclruleId            string
	)

	if v, ok := d.GetOk("nat_gateway_id"); ok {
		request.NatGatewayId = helper.String(v.(string))
		natGatewayId = v.(string)
	}

	if v, ok := d.GetOk("translation_direction"); ok {
		request.TranslationDirection = helper.String(v.(string))
		translationDirection = v.(string)
	}

	if v, ok := d.GetOk("translation_type"); ok {
		request.TranslationType = helper.String(v.(string))
		translationType = v.(string)
	}

	if v, ok := d.GetOk("translation_ip"); ok {
		request.TranslationIp = helper.String(v.(string))
		translationIp = v.(string)
	}

	if v, ok := d.GetOk("translation_acl_rules"); ok {
		for _, item := range v.([]interface{}) {
			translationAclRulesMap := item.(map[string]interface{})
			translationAclRule := vpcv20170312.TranslationAclRule{}
			if v, ok := translationAclRulesMap["protocol"].(string); ok && v != "" {
				translationAclRule.Protocol = helper.String(v)
			}

			if v, ok := translationAclRulesMap["source_port"].(string); ok && v != "" {
				translationAclRule.SourcePort = helper.String(v)
			}

			if v, ok := translationAclRulesMap["source_cidr"].(string); ok && v != "" {
				translationAclRule.SourceCidr = helper.String(v)
			}

			if v, ok := translationAclRulesMap["destination_port"].(string); ok && v != "" {
				translationAclRule.DestinationPort = helper.String(v)
			}

			if v, ok := translationAclRulesMap["destination_cidr"].(string); ok && v != "" {
				translationAclRule.DestinationCidr = helper.String(v)
			}

			if v, ok := translationAclRulesMap["action"].(int); ok {
				translationAclRule.Action = helper.IntUint64(v)
			}

			if v, ok := translationAclRulesMap["description"].(string); ok && v != "" {
				translationAclRule.Description = helper.String(v)
			}

			request.TranslationAclRules = append(request.TranslationAclRules, &translationAclRule)
		}
	}

	if v, ok := d.GetOk("original_ip"); ok {
		request.OriginalIp = helper.String(v.(string))
		originalIp = v.(string)
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreatePrivateNatGatewayTranslationAclRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.TranslationAclRuleSet == nil || len(result.Response.TranslationAclRuleSet) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Create vpc private nat gateway translation acl rule failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create vpc private nat gateway translation acl rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.TranslationAclRuleSet[0].AclRuleId == nil {
		return fmt.Errorf("AclRuleId is nil.")
	}

	aclruleId = helper.UInt64ToStr(*response.Response.TranslationAclRuleSet[0].AclRuleId)

	if originalIp != "" {
		// NETWORK_LAYER
		d.SetId(strings.Join([]string{natGatewayId, translationDirection, translationType, translationIp, originalIp, aclruleId}, tccommon.FILED_SP))
	} else {
		// TRANSPORT_LAYER
		d.SetId(strings.Join([]string{natGatewayId, translationDirection, translationType, translationIp, aclruleId}, tccommon.FILED_SP))
	}

	return resourceTencentCloudVpcPrivateNatGatewayTranslationAclRuleRead(d, meta)
}

func resourceTencentCloudVpcPrivateNatGatewayTranslationAclRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_private_nat_gateway_translation_acl_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                = tccommon.GetLogId(tccommon.ContextNil)
		ctx                  = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service              = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		natGatewayId         string
		translationDirection string
		translationType      string
		translationIp        string
		originalIp           string
		aclruleId            string
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if !(len(idSplit) == 6 || len(idSplit) == 5) {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	if len(idSplit) == 6 {
		natGatewayId = idSplit[0]
		translationDirection = idSplit[1]
		translationType = idSplit[2]
		translationIp = idSplit[3]
		originalIp = idSplit[4]
		aclruleId = idSplit[5]
	} else if len(idSplit) == 5 {
		natGatewayId = idSplit[0]
		translationDirection = idSplit[1]
		translationType = idSplit[2]
		translationIp = idSplit[3]
		aclruleId = idSplit[4]
	} else {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	respData, err := service.DescribeVpcPrivateNatGatewayTranslationAclRuleById(ctx, natGatewayId, translationDirection, translationType, translationIp, originalIp, aclruleId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_vpc_private_nat_gateway_translation_acl_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("nat_gateway_id", natGatewayId)
	_ = d.Set("translation_direction", translationDirection)
	_ = d.Set("translation_type", translationType)
	_ = d.Set("translation_ip", translationIp)
	if originalIp != "" {
		_ = d.Set("original_ip", originalIp)
	}

	translationAclRuleSetList := make([]map[string]interface{}, 0, len(respData))
	for _, translationAclRuleSet := range respData {
		translationAclRuleSetMap := map[string]interface{}{}
		if translationAclRuleSet.Protocol != nil {
			translationAclRuleSetMap["protocol"] = translationAclRuleSet.Protocol
		}

		if translationAclRuleSet.SourcePort != nil {
			translationAclRuleSetMap["source_port"] = translationAclRuleSet.SourcePort
		}

		if translationAclRuleSet.SourceCidr != nil {
			translationAclRuleSetMap["source_cidr"] = translationAclRuleSet.SourceCidr
		}

		if translationAclRuleSet.DestinationPort != nil {
			translationAclRuleSetMap["destination_port"] = translationAclRuleSet.DestinationPort
		}

		if translationAclRuleSet.DestinationCidr != nil {
			translationAclRuleSetMap["destination_cidr"] = translationAclRuleSet.DestinationCidr
		}

		if translationAclRuleSet.AclRuleId != nil {
			translationAclRuleSetMap["acl_rule_id"] = translationAclRuleSet.AclRuleId
		}

		if translationAclRuleSet.Action != nil {
			translationAclRuleSetMap["action"] = translationAclRuleSet.Action
		}

		if translationAclRuleSet.Description != nil {
			translationAclRuleSetMap["description"] = translationAclRuleSet.Description
		}

		translationAclRuleSetList = append(translationAclRuleSetList, translationAclRuleSetMap)
	}

	_ = d.Set("translation_acl_rules", translationAclRuleSetList)

	return nil
}

func resourceTencentCloudVpcPrivateNatGatewayTranslationAclRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_private_nat_gateway_translation_acl_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                = tccommon.GetLogId(tccommon.ContextNil)
		ctx                  = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		natGatewayId         string
		translationDirection string
		translationType      string
		translationIp        string
		originalIp           string
		aclruleId            string
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if !(len(idSplit) == 6 || len(idSplit) == 5) {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	if len(idSplit) == 6 {
		natGatewayId = idSplit[0]
		translationDirection = idSplit[1]
		translationType = idSplit[2]
		translationIp = idSplit[3]
		originalIp = idSplit[4]
		aclruleId = idSplit[5]
	} else if len(idSplit) == 5 {
		natGatewayId = idSplit[0]
		translationDirection = idSplit[1]
		translationType = idSplit[2]
		translationIp = idSplit[3]
		aclruleId = idSplit[4]
	} else {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	if d.HasChange("translation_acl_rules") {
		request := vpcv20170312.NewModifyPrivateNatGatewayTranslationAclRuleRequest()
		if v, ok := d.GetOk("translation_acl_rules"); ok {
			for _, item := range v.([]interface{}) {
				translationAclRulesMap := item.(map[string]interface{})
				translationAclRule := vpcv20170312.TranslationAclRule{}
				if v, ok := translationAclRulesMap["protocol"].(string); ok && v != "" {
					translationAclRule.Protocol = helper.String(v)
				}

				if v, ok := translationAclRulesMap["source_port"].(string); ok && v != "" {
					translationAclRule.SourcePort = helper.String(v)
				}

				if v, ok := translationAclRulesMap["source_cidr"].(string); ok && v != "" {
					translationAclRule.SourceCidr = helper.String(v)
				}

				if v, ok := translationAclRulesMap["destination_port"].(string); ok && v != "" {
					translationAclRule.DestinationPort = helper.String(v)
				}

				if v, ok := translationAclRulesMap["destination_cidr"].(string); ok && v != "" {
					translationAclRule.DestinationCidr = helper.String(v)
				}

				translationAclRule.AclRuleId = helper.StrToUint64Point(aclruleId)

				if v, ok := translationAclRulesMap["action"].(int); ok {
					translationAclRule.Action = helper.IntUint64(v)
				}

				if v, ok := translationAclRulesMap["description"].(string); ok && v != "" {
					translationAclRule.Description = helper.String(v)
				}

				request.TranslationAclRules = append(request.TranslationAclRules, &translationAclRule)
			}
		}

		request.NatGatewayId = &natGatewayId
		request.TranslationDirection = &translationDirection
		request.TranslationType = &translationType
		request.TranslationIp = &translationIp
		if originalIp != "" {
			request.OriginalIp = &originalIp
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyPrivateNatGatewayTranslationAclRuleWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update vpc private nat gateway translation acl rule failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudVpcPrivateNatGatewayTranslationAclRuleRead(d, meta)
}

func resourceTencentCloudVpcPrivateNatGatewayTranslationAclRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_private_nat_gateway_translation_acl_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                = tccommon.GetLogId(tccommon.ContextNil)
		ctx                  = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request              = vpcv20170312.NewDeletePrivateNatGatewayTranslationAclRuleRequest()
		natGatewayId         string
		translationDirection string
		translationType      string
		translationIp        string
		originalIp           string
		aclruleId            string
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if !(len(idSplit) == 6 || len(idSplit) == 5) {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	if len(idSplit) == 6 {
		natGatewayId = idSplit[0]
		translationDirection = idSplit[1]
		translationType = idSplit[2]
		translationIp = idSplit[3]
		originalIp = idSplit[4]
		aclruleId = idSplit[5]
	} else if len(idSplit) == 5 {
		natGatewayId = idSplit[0]
		translationDirection = idSplit[1]
		translationType = idSplit[2]
		translationIp = idSplit[3]
		aclruleId = idSplit[4]
	} else {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	request.NatGatewayId = &natGatewayId
	request.TranslationDirection = &translationDirection
	request.TranslationType = &translationType
	request.TranslationIp = &translationIp
	request.AclRuleIds = append(request.AclRuleIds, helper.StrToUint64Point(aclruleId))
	if originalIp != "" {
		request.OriginalIp = &originalIp
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DeletePrivateNatGatewayTranslationAclRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete vpc private nat gateway translation acl rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
