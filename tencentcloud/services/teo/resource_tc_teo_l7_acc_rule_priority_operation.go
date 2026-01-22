package teo

import (
	"log"

	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudTeoL7AccRulePriorityOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoL7AccRulePriorityOperationCreate,
		Read:   resourceTencentCloudTeoL7AccRulePriorityOperationRead,
		Delete: resourceTencentCloudTeoL7AccRulePriorityOperationDelete,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone id.",
			},
			"rule_ids": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Complete list of rule IDs under site ID.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTencentCloudTeoL7AccRulePriorityOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l7_acc_rule_priority_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := teov20220901.NewModifyL7AccRulePriorityRequest()
	zoneId := d.Get("zone_id").(string)
	request.ZoneId = helper.String(zoneId)
	ruleIds := d.Get("rule_ids").([]interface{})
	for _, ruleId := range ruleIds {
		request.RuleIds = append(request.RuleIds, helper.String(ruleId.(string)))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyL7AccRulePriority(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(zoneId)

	return resourceTencentCloudTeoL7AccRulePriorityOperationRead(d, meta)
}

func resourceTencentCloudTeoL7AccRulePriorityOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l7_acc_rule_priority_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTeoL7AccRulePriorityOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l7_acc_rule_priority_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
