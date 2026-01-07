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
		Schema: map[string]*schema.Schema{
			"uuid_list": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The unique IDs of the rule, which is not required when you create a rule. The priority will be determined by the index position of the UUID in the list.",
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

	if v, ok := d.GetOk("uuid_list"); ok {
		for _, item := range v.([]interface{}) {
			uuid := helper.IntToStr(item.(int))
			natPolicy, err := service.DescribeCfwNatPolicyOrderConfigById(ctx, uuid)
			if err != nil {
				return err
			}

			if natPolicy == nil {
				return fmt.Errorf("uuid %d does not exist.", item.(int))
			}

			if natPolicy.Uuid == nil {
				return fmt.Errorf("uuid %d does not exist.", item.(int))
			}
		}
	}

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

	if d.HasChange("uuid_list") {
		if v, ok := d.GetOk("uuid_list"); ok {
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
				if natPolicy.SourceContent != nil {
					modifyRuleItem.SourceContent = natPolicy.SourceContent
				}

				if natPolicy.SourceType != nil {
					modifyRuleItem.SourceType = natPolicy.SourceType
				}

				if natPolicy.TargetContent != nil {
					modifyRuleItem.TargetContent = natPolicy.TargetContent
				}

				if natPolicy.TargetType != nil {
					modifyRuleItem.TargetType = natPolicy.TargetType
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
