package cfw

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfwv20190904 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCfwNatPolicyOrderConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfwNatPolicyOrderConfigCreate,
		Read:   resourceTencentCloudCfwNatPolicyOrderConfigRead,
		Update: resourceTencentCloudCfwNatPolicyOrderConfigUpdate,
		Delete: resourceTencentCloudCfwNatPolicyOrderConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"inbound_rule_uuid_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The unique IDs of the inbound rule, which is not required when you create a rule. The priority will be determined by the index position of the UUID in the list.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"outbound_rule_uuid_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The unique IDs of the outbound rule, which is not required when you create a rule. The priority will be determined by the index position of the UUID in the list.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func resourceTencentCloudCfwNatPolicyOrderConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_nat_policy_order_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	d.SetId(helper.BuildToken())

	return resourceTencentCloudCfwNatPolicyOrderConfigUpdate(d, meta)
}

func resourceTencentCloudCfwNatPolicyOrderConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_nat_policy_order_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	respData, err := service.DescribeCfwNatPolicyOrderConfigs(ctx)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_cfw_nat_policy_order_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	inTmpList := make([]int, 0, len(respData))
	outTmpList := make([]int, 0, len(respData))
	for _, item := range respData {
		if item != nil && item.Uuid != nil && item.Direction != nil {
			if *item.Direction == 1 {
				inTmpList = append(inTmpList, int(*item.Uuid))
			} else if *item.Direction == 0 {
				outTmpList = append(outTmpList, int(*item.Uuid))
			}
		}
	}

	_ = d.Set("inbound_rule_uuid_list", inTmpList)
	_ = d.Set("outbound_rule_uuid_list", outTmpList)

	return nil
}

func resourceTencentCloudCfwNatPolicyOrderConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_nat_policy_order_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	if d.HasChange("inbound_rule_uuid_list") {
		if v, ok := d.GetOk("inbound_rule_uuid_list"); ok {
			for k, item := range v.([]interface{}) {
				request := cfwv20190904.NewModifyNatAcRuleRequest()
				uuid := helper.IntToStr(item.(int))
				natPolicy, err := service.DescribeCfwNatPolicyOrderConfigById(ctx, uuid)
				if err != nil {
					return err
				}

				if natPolicy == nil {
					return fmt.Errorf("uuid %d does not exist.", item.(int))
				}

				modifyRuleItem := cfwv20190904.CreateNatRuleItem{}
				if natPolicy.SourceType != nil {
					modifyRuleItem.SourceType = natPolicy.SourceType
				}

				if natPolicy.SourceContent != nil {
					if natPolicy.SourceType != nil && *natPolicy.SourceType == "tag" {
						ref, _ := service.SplitAndFormat(*natPolicy.SourceContent)
						modifyRuleItem.SourceContent = &ref
					} else {
						modifyRuleItem.SourceContent = natPolicy.SourceContent
					}
				}

				if natPolicy.TargetType != nil {
					modifyRuleItem.TargetType = natPolicy.TargetType
				}

				if natPolicy.TargetContent != nil {
					if natPolicy.TargetType != nil && *natPolicy.TargetType == "tag" {
						ref, _ := service.SplitAndFormat(*natPolicy.TargetContent)
						modifyRuleItem.TargetContent = &ref
					} else {
						modifyRuleItem.TargetContent = natPolicy.TargetContent
					}
				}

				if natPolicy.Protocol != nil {
					modifyRuleItem.Protocol = natPolicy.Protocol
				}

				if natPolicy.RuleAction != nil {
					modifyRuleItem.RuleAction = natPolicy.RuleAction
				}

				if natPolicy.Port != nil {
					modifyRuleItem.Port = natPolicy.Port
				}

				if natPolicy.Direction != nil {
					modifyRuleItem.Direction = natPolicy.Direction
				}

				if natPolicy.Enable != nil {
					modifyRuleItem.Enable = natPolicy.Enable
				}

				if natPolicy.Description != nil {
					modifyRuleItem.Description = natPolicy.Description
				}

				if natPolicy.ParamTemplateId != nil {
					modifyRuleItem.ParamTemplateId = natPolicy.ParamTemplateId
				}

				if natPolicy.Scope != nil {
					modifyRuleItem.Scope = natPolicy.Scope
				}

				orderIndex := k + 1
				modifyRuleItem.OrderIndex = helper.IntInt64(orderIndex)
				modifyRuleItem.Uuid = helper.IntInt64(item.(int))
				request.Rules = append(request.Rules, &modifyRuleItem)
				reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwV20190904Client().ModifyNatAcRuleWithContext(ctx, request)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					return nil
				})

				if reqErr != nil {
					log.Printf("[CRITAL]%s update cfw nat policy order config failed, reason:%+v", logId, reqErr)
					return reqErr
				}
			}
		}
	}

	if d.HasChange("outbound_rule_uuid_list") {
		if v, ok := d.GetOk("outbound_rule_uuid_list"); ok {
			for k, item := range v.([]interface{}) {
				request := cfwv20190904.NewModifyNatAcRuleRequest()
				uuid := helper.IntToStr(item.(int))
				natPolicy, err := service.DescribeCfwNatPolicyOrderConfigById(ctx, uuid)
				if err != nil {
					return err
				}

				if natPolicy == nil {
					return fmt.Errorf("uuid %d does not exist.", item.(int))
				}

				modifyRuleItem := cfwv20190904.CreateNatRuleItem{}
				if natPolicy.SourceType != nil {
					modifyRuleItem.SourceType = natPolicy.SourceType
				}

				if natPolicy.SourceContent != nil {
					if natPolicy.SourceType != nil && *natPolicy.SourceType == "tag" {
						ref, _ := service.SplitAndFormat(*natPolicy.SourceContent)
						modifyRuleItem.SourceContent = &ref
					} else {
						modifyRuleItem.SourceContent = natPolicy.SourceContent
					}
				}

				if natPolicy.TargetType != nil {
					modifyRuleItem.TargetType = natPolicy.TargetType
				}

				if natPolicy.TargetContent != nil {
					if natPolicy.TargetType != nil && *natPolicy.TargetType == "tag" {
						ref, _ := service.SplitAndFormat(*natPolicy.TargetContent)
						modifyRuleItem.TargetContent = &ref
					} else {
						modifyRuleItem.TargetContent = natPolicy.TargetContent
					}
				}

				if natPolicy.Protocol != nil {
					modifyRuleItem.Protocol = natPolicy.Protocol
				}

				if natPolicy.RuleAction != nil {
					modifyRuleItem.RuleAction = natPolicy.RuleAction
				}

				if natPolicy.Port != nil {
					modifyRuleItem.Port = natPolicy.Port
				}

				if natPolicy.Direction != nil {
					modifyRuleItem.Direction = natPolicy.Direction
				}

				if natPolicy.Enable != nil {
					modifyRuleItem.Enable = natPolicy.Enable
				}

				if natPolicy.Description != nil {
					modifyRuleItem.Description = natPolicy.Description
				}

				if natPolicy.ParamTemplateId != nil {
					modifyRuleItem.ParamTemplateId = natPolicy.ParamTemplateId
				}

				if natPolicy.Scope != nil {
					modifyRuleItem.Scope = natPolicy.Scope
				}

				orderIndex := k + 1
				modifyRuleItem.OrderIndex = helper.IntInt64(orderIndex)
				modifyRuleItem.Uuid = helper.IntInt64(item.(int))
				request.Rules = append(request.Rules, &modifyRuleItem)
				reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwV20190904Client().ModifyNatAcRuleWithContext(ctx, request)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					return nil
				})

				if reqErr != nil {
					log.Printf("[CRITAL]%s update cfw nat policy order config failed, reason:%+v", logId, reqErr)
					return reqErr
				}
			}
		}
	}

	return resourceTencentCloudCfwNatPolicyOrderConfigRead(d, meta)
}

func resourceTencentCloudCfwNatPolicyOrderConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_nat_policy_order_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
