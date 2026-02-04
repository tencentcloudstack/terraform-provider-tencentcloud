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

func ResourceTencentCloudCfwVpcPolicyOrderConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfwVpcPolicyOrderConfigCreate,
		Read:   resourceTencentCloudCfwVpcPolicyOrderConfigRead,
		Update: resourceTencentCloudCfwVpcPolicyOrderConfigUpdate,
		Delete: resourceTencentCloudCfwVpcPolicyOrderConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"rule_uuid_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The unique IDs of the rule, which is not required when you create a rule. The priority will be determined by the index position of the UUID in the list.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func resourceTencentCloudCfwVpcPolicyOrderConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_vpc_policy_order_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	d.SetId(helper.BuildToken())

	return resourceTencentCloudCfwVpcPolicyOrderConfigUpdate(d, meta)
}

func resourceTencentCloudCfwVpcPolicyOrderConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_vpc_policy_order_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	respData, err := service.DescribeCfwVpcPolicyOrderConfigs(ctx)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_cfw_vpc_policy_order_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	tmpList := make([]int, 0, len(respData))
	for _, item := range respData {
		if item != nil && item.Uuid != nil {
			tmpList = append(tmpList, int(*item.Uuid))
		}
	}

	_ = d.Set("rule_uuid_list", tmpList)

	return nil
}

func resourceTencentCloudCfwVpcPolicyOrderConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_vpc_policy_order_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	if d.HasChange("rule_uuid_list") {
		if v, ok := d.GetOk("rule_uuid_list"); ok {
			for k, item := range v.([]interface{}) {
				request := cfwv20190904.NewModifyVpcAcRuleRequest()
				uuid := helper.IntToStr(item.(int))
				natPolicy, err := service.DescribeCfwVpcPolicyOrderConfigById(ctx, uuid)
				if err != nil {
					return err
				}

				if natPolicy == nil {
					return fmt.Errorf("uuid %d does not exist.", item.(int))
				}

				modifyRuleItem := cfwv20190904.VpcRuleItem{}
				if natPolicy.SourceContent != nil {
					modifyRuleItem.SourceContent = natPolicy.SourceContent
				}

				if natPolicy.SourceType != nil {
					modifyRuleItem.SourceType = natPolicy.SourceType
				}

				if natPolicy.DestContent != nil {
					modifyRuleItem.DestContent = natPolicy.DestContent
				}

				if natPolicy.DestType != nil {
					modifyRuleItem.DestType = natPolicy.DestType
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

				if natPolicy.Description != nil {
					modifyRuleItem.Description = natPolicy.Description
				}

				if natPolicy.Enable != nil {
					modifyRuleItem.Enable = natPolicy.Enable
				}

				if natPolicy.FwGroupId != nil {
					modifyRuleItem.FwGroupId = natPolicy.FwGroupId
				}

				orderIndex := k + 1
				modifyRuleItem.OrderIndex = helper.IntInt64(orderIndex)
				modifyRuleItem.Uuid = helper.IntInt64(item.(int))
				request.Rules = append(request.Rules, &modifyRuleItem)
				reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwV20190904Client().ModifyVpcAcRuleWithContext(ctx, request)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					return nil
				})

				if reqErr != nil {
					log.Printf("[CRITAL]%s update cfw vpc policy order config failed, reason:%+v", logId, reqErr)
					return reqErr
				}
			}
		}
	}

	return resourceTencentCloudCfwVpcPolicyOrderConfigRead(d, meta)
}

func resourceTencentCloudCfwVpcPolicyOrderConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_vpc_policy_order_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
