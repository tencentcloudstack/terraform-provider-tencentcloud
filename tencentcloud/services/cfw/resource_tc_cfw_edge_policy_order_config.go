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

func ResourceTencentCloudCfwEdgePolicyOrderConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfwEdgePolicyOrderConfigCreate,
		Read:   resourceTencentCloudCfwEdgePolicyOrderConfigRead,
		Update: resourceTencentCloudCfwEdgePolicyOrderConfigUpdate,
		Delete: resourceTencentCloudCfwEdgePolicyOrderConfigDelete,
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

func resourceTencentCloudCfwEdgePolicyOrderConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_edge_policy_order_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	d.SetId(helper.BuildToken())

	return resourceTencentCloudCfwEdgePolicyOrderConfigUpdate(d, meta)
}

func resourceTencentCloudCfwEdgePolicyOrderConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_edge_policy_order_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	respData, err := service.DescribeCfwEdgePolicyOrderConfigs(ctx)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_cfw_edge_policy_order_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
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

func resourceTencentCloudCfwEdgePolicyOrderConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_edge_policy_order_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	if d.HasChange("inbound_rule_uuid_list") {
		if v, ok := d.GetOk("inbound_rule_uuid_list"); ok {
			for k, item := range v.([]interface{}) {
				request := cfwv20190904.NewModifyAclRuleRequest()
				uuid := helper.IntToStr(item.(int))
				edgePolicy, err := service.DescribeCfwEdgePolicyOrderConfigById(ctx, uuid)
				if err != nil {
					return err
				}

				if edgePolicy == nil {
					return fmt.Errorf("uuid %d does not exist.", item.(int))
				}

				modifyRuleItem := cfwv20190904.CreateRuleItem{}
				if edgePolicy.SourceType != nil {
					modifyRuleItem.SourceType = edgePolicy.SourceType
				}

				if edgePolicy.SourceContent != nil {
					if edgePolicy.SourceType != nil && *edgePolicy.SourceType == "tag" {
						ref, _ := service.SplitAndFormat(*edgePolicy.SourceContent)
						modifyRuleItem.SourceContent = &ref
					} else {
						modifyRuleItem.SourceContent = edgePolicy.SourceContent
					}
				}

				if edgePolicy.TargetType != nil {
					modifyRuleItem.TargetType = edgePolicy.TargetType
				}

				if edgePolicy.TargetContent != nil {
					if edgePolicy.TargetType != nil && *edgePolicy.TargetType == "tag" {
						ref, _ := service.SplitAndFormat(*edgePolicy.TargetContent)
						modifyRuleItem.TargetContent = &ref
					} else {
						modifyRuleItem.TargetContent = edgePolicy.TargetContent
					}
				}

				if edgePolicy.Protocol != nil {
					modifyRuleItem.Protocol = edgePolicy.Protocol
				}

				if edgePolicy.RuleAction != nil {
					modifyRuleItem.RuleAction = edgePolicy.RuleAction
				}

				if edgePolicy.Port != nil {
					modifyRuleItem.Port = edgePolicy.Port
				}

				if edgePolicy.Direction != nil {
					modifyRuleItem.Direction = edgePolicy.Direction
				}

				if edgePolicy.Enable != nil {
					modifyRuleItem.Enable = edgePolicy.Enable
				}

				if edgePolicy.Description != nil {
					modifyRuleItem.Description = edgePolicy.Description
				}

				if edgePolicy.ParamTemplateId != nil {
					modifyRuleItem.ParamTemplateId = edgePolicy.ParamTemplateId
				}

				if edgePolicy.Scope != nil {
					modifyRuleItem.Scope = edgePolicy.Scope
				}

				orderIndex := k + 1
				modifyRuleItem.OrderIndex = helper.IntInt64(orderIndex)
				modifyRuleItem.Uuid = helper.IntInt64(item.(int))
				request.Rules = append(request.Rules, &modifyRuleItem)
				reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwV20190904Client().ModifyAclRuleWithContext(ctx, request)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					return nil
				})

				if reqErr != nil {
					log.Printf("[CRITAL]%s update cfw edge policy order config failed, reason:%+v", logId, reqErr)
					return reqErr
				}
			}
		}
	}

	if d.HasChange("outbound_rule_uuid_list") {
		if v, ok := d.GetOk("outbound_rule_uuid_list"); ok {
			for k, item := range v.([]interface{}) {
				request := cfwv20190904.NewModifyAclRuleRequest()
				uuid := helper.IntToStr(item.(int))
				edgePolicy, err := service.DescribeCfwEdgePolicyOrderConfigById(ctx, uuid)
				if err != nil {
					return err
				}

				if edgePolicy == nil {
					return fmt.Errorf("uuid %d does not exist.", item.(int))
				}

				modifyRuleItem := cfwv20190904.CreateRuleItem{}
				if edgePolicy.SourceType != nil {
					modifyRuleItem.SourceType = edgePolicy.SourceType
				}

				if edgePolicy.SourceContent != nil {
					if edgePolicy.SourceType != nil && *edgePolicy.SourceType == "tag" {
						ref, _ := service.SplitAndFormat(*edgePolicy.SourceContent)
						modifyRuleItem.SourceContent = &ref
					} else {
						modifyRuleItem.SourceContent = edgePolicy.SourceContent
					}
				}

				if edgePolicy.TargetType != nil {
					modifyRuleItem.TargetType = edgePolicy.TargetType
				}

				if edgePolicy.TargetContent != nil {
					if edgePolicy.TargetType != nil && *edgePolicy.TargetType == "tag" {
						ref, _ := service.SplitAndFormat(*edgePolicy.TargetContent)
						modifyRuleItem.TargetContent = &ref
					} else {
						modifyRuleItem.TargetContent = edgePolicy.TargetContent
					}
				}

				if edgePolicy.Protocol != nil {
					modifyRuleItem.Protocol = edgePolicy.Protocol
				}

				if edgePolicy.RuleAction != nil {
					modifyRuleItem.RuleAction = edgePolicy.RuleAction
				}

				if edgePolicy.Port != nil {
					modifyRuleItem.Port = edgePolicy.Port
				}

				if edgePolicy.Direction != nil {
					modifyRuleItem.Direction = edgePolicy.Direction
				}

				if edgePolicy.Enable != nil {
					modifyRuleItem.Enable = edgePolicy.Enable
				}

				if edgePolicy.Description != nil {
					modifyRuleItem.Description = edgePolicy.Description
				}

				if edgePolicy.ParamTemplateId != nil {
					modifyRuleItem.ParamTemplateId = edgePolicy.ParamTemplateId
				}

				if edgePolicy.Scope != nil {
					modifyRuleItem.Scope = edgePolicy.Scope
				}

				orderIndex := k + 1
				modifyRuleItem.OrderIndex = helper.IntInt64(orderIndex)
				modifyRuleItem.Uuid = helper.IntInt64(item.(int))
				request.Rules = append(request.Rules, &modifyRuleItem)
				reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwV20190904Client().ModifyAclRuleWithContext(ctx, request)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					return nil
				})

				if reqErr != nil {
					log.Printf("[CRITAL]%s update cfw edge policy order config failed, reason:%+v", logId, reqErr)
					return reqErr
				}
			}
		}
	}

	return resourceTencentCloudCfwEdgePolicyOrderConfigRead(d, meta)
}

func resourceTencentCloudCfwEdgePolicyOrderConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_edge_policy_order_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
