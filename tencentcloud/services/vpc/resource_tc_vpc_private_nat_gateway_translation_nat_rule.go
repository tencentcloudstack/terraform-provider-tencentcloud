package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpcv20170312 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

var MAX_CREATE_RULES_LEN = 20

func ResourceTencentCloudVpcPrivateNatGatewayTranslationNatRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcPrivateNatGatewayTranslationNatRuleCreate,
		Read:   resourceTencentCloudVpcPrivateNatGatewayTranslationNatRuleRead,
		Update: resourceTencentCloudVpcPrivateNatGatewayTranslationNatRuleUpdate,
		Delete: resourceTencentCloudVpcPrivateNatGatewayTranslationNatRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"nat_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Private NAT gateway unique ID, such as: `intranat-xxxxxxxx`.",
			},

			"translation_nat_rules": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Translation rule object array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"translation_direction": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Translation rule target, optional values \"LOCAL\",\"PEER\".",
						},
						"translation_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Translation rule type, optional values \"NETWORK_LAYER\",\"TRANSPORT_LAYER\".",
						},
						"translation_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Translation IP, when translation rule type is transport layer, it is an IP pool.",
						},
						"description": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Translation rule description.",
						},
						"original_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Source IP, valid when translation rule type is network layer.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudVpcPrivateNatGatewayTranslationNatRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_private_nat_gateway_translation_nat_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request      = vpcv20170312.NewCreatePrivateNatGatewayTranslationNatRuleRequest()
		natGatewayId string
	)

	if v, ok := d.GetOk("nat_gateway_id"); ok {
		request.NatGatewayId = helper.String(v.(string))
		natGatewayId = v.(string)
	}

	var allRules []*vpcv20170312.TranslationNatRuleInput
	if v, ok := d.GetOk("translation_nat_rules"); ok {
		for _, item := range v.(*schema.Set).List() {
			translationNatRulesMap := item.(map[string]interface{})
			translationNatRuleInput := vpcv20170312.TranslationNatRuleInput{}
			if v, ok := translationNatRulesMap["translation_direction"].(string); ok && v != "" {
				translationNatRuleInput.TranslationDirection = helper.String(v)
			}

			if v, ok := translationNatRulesMap["translation_type"].(string); ok && v != "" {
				translationNatRuleInput.TranslationType = helper.String(v)
			}

			if v, ok := translationNatRulesMap["translation_ip"].(string); ok && v != "" {
				translationNatRuleInput.TranslationIp = helper.String(v)
			}

			if v, ok := translationNatRulesMap["description"].(string); ok && v != "" {
				translationNatRuleInput.Description = helper.String(v)
			}

			if v, ok := translationNatRulesMap["original_ip"].(string); ok && v != "" {
				translationNatRuleInput.OriginalIp = helper.String(v)
			}

			allRules = append(allRules, &translationNatRuleInput)
		}
	}

	for i := 0; i < len(allRules); i += MAX_CREATE_RULES_LEN {
		end := i + MAX_CREATE_RULES_LEN
		if end > len(allRules) {
			end = len(allRules)
		}

		batchRules := allRules[i:end]
		request.TranslationNatRules = batchRules
		request.NatGatewayId = helper.String(natGatewayId)
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreatePrivateNatGatewayTranslationNatRuleWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.NatGatewayId == nil {
				return resource.NonRetryableError(fmt.Errorf("Create vpc private nat gateway translation nat rule failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s create vpc private nat gateway translation nat rule failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	d.SetId(natGatewayId)
	return resourceTencentCloudVpcPrivateNatGatewayTranslationNatRuleRead(d, meta)
}

func resourceTencentCloudVpcPrivateNatGatewayTranslationNatRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_private_nat_gateway_translation_nat_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service      = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		natGatewayId = d.Id()
	)

	respData, err := service.DescribeVpcPrivateNatGatewayTranslationNatRuleById(ctx, natGatewayId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_vpc_private_nat_gateway_translation_nat_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("nat_gateway_id", natGatewayId)

	translationNatRuleSetList := make([]map[string]interface{}, 0, len(respData))
	for _, translationNatRuleSet := range respData {
		translationNatRuleSetMap := map[string]interface{}{}
		if translationNatRuleSet.TranslationDirection != nil {
			translationNatRuleSetMap["translation_direction"] = translationNatRuleSet.TranslationDirection
		}

		if translationNatRuleSet.TranslationType != nil {
			translationNatRuleSetMap["translation_type"] = translationNatRuleSet.TranslationType
		}

		if translationNatRuleSet.TranslationIp != nil {
			translationNatRuleSetMap["translation_ip"] = translationNatRuleSet.TranslationIp
		}

		if translationNatRuleSet.Description != nil {
			translationNatRuleSetMap["description"] = translationNatRuleSet.Description
		}

		if translationNatRuleSet.OriginalIp != nil {
			translationNatRuleSetMap["original_ip"] = translationNatRuleSet.OriginalIp
		}

		translationNatRuleSetList = append(translationNatRuleSetList, translationNatRuleSetMap)
	}

	_ = d.Set("translation_nat_rules", translationNatRuleSetList)

	return nil
}

func resourceTencentCloudVpcPrivateNatGatewayTranslationNatRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_private_nat_gateway_translation_nat_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		natGatewayId = d.Id()
	)

	if d.HasChange("translation_nat_rules") {
		oldInterface, newInterface := d.GetChange("translation_nat_rules")
		olds := oldInterface.(*schema.Set)
		news := newInterface.(*schema.Set)
		remove := olds.Difference(news).List()
		add := news.Difference(olds).List()

		if len(remove) > 0 {
			request := vpcv20170312.NewDeletePrivateNatGatewayTranslationNatRuleRequest()
			for _, item := range remove {
				translationNatRulesMap := item.(map[string]interface{})
				translationNatRuleInput := vpcv20170312.TranslationNatRule{}
				if v, ok := translationNatRulesMap["translation_direction"].(string); ok && v != "" {
					translationNatRuleInput.TranslationDirection = helper.String(v)
				}

				if v, ok := translationNatRulesMap["translation_type"].(string); ok && v != "" {
					translationNatRuleInput.TranslationType = helper.String(v)
				}

				if v, ok := translationNatRulesMap["translation_ip"].(string); ok && v != "" {
					translationNatRuleInput.TranslationIp = helper.String(v)
				}

				if v, ok := translationNatRulesMap["description"].(string); ok && v != "" {
					translationNatRuleInput.Description = helper.String(v)
				}

				if v, ok := translationNatRulesMap["original_ip"].(string); ok && v != "" {
					translationNatRuleInput.OriginalIp = helper.String(v)
				}

				request.TranslationNatRules = append(request.TranslationNatRules, &translationNatRuleInput)
			}

			request.NatGatewayId = &natGatewayId
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DeletePrivateNatGatewayTranslationNatRuleWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s delete vpc private nat gateway translation nat rule failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		}

		if len(add) > 0 {
			request := vpcv20170312.NewCreatePrivateNatGatewayTranslationNatRuleRequest()
			var allRules []*vpcv20170312.TranslationNatRuleInput
			for _, item := range add {
				translationNatRulesMap := item.(map[string]interface{})
				translationNatRuleInput := vpcv20170312.TranslationNatRuleInput{}
				if v, ok := translationNatRulesMap["translation_direction"].(string); ok && v != "" {
					translationNatRuleInput.TranslationDirection = helper.String(v)
				}

				if v, ok := translationNatRulesMap["translation_type"].(string); ok && v != "" {
					translationNatRuleInput.TranslationType = helper.String(v)
				}

				if v, ok := translationNatRulesMap["translation_ip"].(string); ok && v != "" {
					translationNatRuleInput.TranslationIp = helper.String(v)
				}

				if v, ok := translationNatRulesMap["description"].(string); ok && v != "" {
					translationNatRuleInput.Description = helper.String(v)
				}

				if v, ok := translationNatRulesMap["original_ip"].(string); ok && v != "" {
					translationNatRuleInput.OriginalIp = helper.String(v)
				}

				allRules = append(allRules, &translationNatRuleInput)
			}

			for i := 0; i < len(allRules); i += MAX_CREATE_RULES_LEN {
				end := i + MAX_CREATE_RULES_LEN
				if end > len(allRules) {
					end = len(allRules)
				}

				batchRules := allRules[i:end]
				request.TranslationNatRules = batchRules
				request.NatGatewayId = helper.String(natGatewayId)
				reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreatePrivateNatGatewayTranslationNatRuleWithContext(ctx, request)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					if result == nil || result.Response == nil || result.Response.NatGatewayId == nil {
						return resource.NonRetryableError(fmt.Errorf("Create vpc private nat gateway translation nat rule failed, Response is nil."))
					}

					return nil
				})

				if reqErr != nil {
					log.Printf("[CRITAL]%s create vpc private nat gateway translation nat rule failed, reason:%+v", logId, reqErr)
					return reqErr
				}
			}
		}
	}

	return resourceTencentCloudVpcPrivateNatGatewayTranslationNatRuleRead(d, meta)
}

func resourceTencentCloudVpcPrivateNatGatewayTranslationNatRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_private_nat_gateway_translation_nat_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service      = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request      = vpcv20170312.NewDeletePrivateNatGatewayTranslationNatRuleRequest()
		natGatewayId = d.Id()
	)

	// get all rules
	respData, err := service.DescribeVpcPrivateNatGatewayTranslationNatRuleById(ctx, natGatewayId)
	if err != nil {
		return err
	}

	for _, item := range respData {
		translationNatRule := vpcv20170312.TranslationNatRule{}
		if item.TranslationDirection != nil {
			translationNatRule.TranslationDirection = item.TranslationDirection
		}

		if item.TranslationType != nil {
			translationNatRule.TranslationType = item.TranslationType
		}

		if item.TranslationIp != nil {
			translationNatRule.TranslationIp = item.TranslationIp
		}

		if item.Description != nil {
			translationNatRule.Description = item.Description
		}

		if item.OriginalIp != nil {
			translationNatRule.OriginalIp = item.OriginalIp
		}

		request.TranslationNatRules = append(request.TranslationNatRules, &translationNatRule)
	}

	request.NatGatewayId = &natGatewayId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DeletePrivateNatGatewayTranslationNatRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete vpc private nat gateway translation nat rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
