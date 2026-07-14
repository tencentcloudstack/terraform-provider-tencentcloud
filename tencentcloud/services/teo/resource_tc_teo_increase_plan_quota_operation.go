package teo

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoIncreasePlanQuotaOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoIncreasePlanQuotaOperationCreate,
		Read:   resourceTencentCloudTeoIncreasePlanQuotaOperationRead,
		Delete: resourceTencentCloudTeoIncreasePlanQuotaOperationDelete,
		Schema: map[string]*schema.Schema{
			"plan_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Plan ID, e.g., edgeone-2unuvzjmmn2q.",
			},
			"quota_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Quota type. Valid values: `site` (site count), `precise_access_control_rule` (Web Protection - Custom Rules - Precise Match Policy rule quota), `rate_limiting_rule` (Web Protection - Rate Limiting - Precise Rate Limiting module rule quota).",
			},
			"quota_number": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Number of quotas to increase. Maximum is 100 per request.",
			},
			"deal_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Order number returned after successful quota increase.",
			},
		},
	}
}

func resourceTencentCloudTeoIncreasePlanQuotaOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_increase_plan_quota.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewIncreasePlanQuotaRequest()
	)

	planId := d.Get("plan_id").(string)
	request.PlanId = helper.String(planId)

	quotaType := d.Get("quota_type").(string)
	request.QuotaType = helper.String(quotaType)

	quotaNumber := int64(d.Get("quota_number").(int))
	request.QuotaNumber = &quotaNumber

	var dealName string
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().IncreasePlanQuotaWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result == nil || result.Response == nil {
			return tccommon.RetryError(nil)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		if result.Response.DealName != nil {
			dealName = *result.Response.DealName
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s increase plan quota failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(helper.BuildToken())

	if dealName != "" {
		_ = d.Set("deal_name", dealName)
	}

	return resourceTencentCloudTeoIncreasePlanQuotaOperationRead(d, meta)
}

func resourceTencentCloudTeoIncreasePlanQuotaOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_increase_plan_quota.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTeoIncreasePlanQuotaOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_increase_plan_quota.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
