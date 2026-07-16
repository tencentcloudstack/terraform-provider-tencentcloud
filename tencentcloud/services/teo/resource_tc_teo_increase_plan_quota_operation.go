package teo

import (
	"fmt"
	"log"

	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Description: "Plan ID, in the format of edgeone-xxxxxxxx.",
			},
			"quota_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Quota type to increase. Valid values: site, precise_access_control_rule, rate_limiting_rule.",
			},
			"quota_number": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Number of quotas to increase. Maximum 100 per request.",
			},
			"deal_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Order number returned by the API.",
			},
		},
	}
}

func resourceTencentCloudTeoIncreasePlanQuotaOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_increase_plan_quota_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	planId := d.Get("plan_id").(string)
	quotaType := d.Get("quota_type").(string)
	quotaNumber := int64(d.Get("quota_number").(int))

	if planId == "" {
		return fmt.Errorf("plan_id is required")
	}
	if quotaType == "" {
		return fmt.Errorf("quota_type is required")
	}

	request := teov20220901.NewIncreasePlanQuotaRequest()
	request.PlanId = helper.String(planId)
	request.QuotaType = helper.String(quotaType)
	request.QuotaNumber = helper.Int64(quotaNumber)

	var dealName string
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().IncreasePlanQuota(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.DealName == nil {
			return resource.NonRetryableError(fmt.Errorf("IncreasePlanQuota API returned empty response"))
		}
		dealName = *result.Response.DealName
		return nil
	})
	if err != nil {
		return err
	}

	log.Printf("[DEBUG]%s api[%s] success, deal_name: %s\n", logId, request.GetAction(), dealName)

	d.SetId(helper.BuildToken())
	_ = d.Set("deal_name", dealName)

	return resourceTencentCloudTeoIncreasePlanQuotaOperationRead(d, meta)
}

func resourceTencentCloudTeoIncreasePlanQuotaOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_increase_plan_quota_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTeoIncreasePlanQuotaOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_increase_plan_quota_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
