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

func ResourceTencentCloudTeoBindZoneToPlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoBindZoneToPlanCreate,
		Read:   resourceTencentCloudTeoBindZoneToPlanRead,
		Delete: resourceTencentCloudTeoBindZoneToPlanDelete,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the unbound zone that needs to be bound to a plan.",
			},
			"plan_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the target plan to be bound.",
			},
		},
	}
}

func resourceTencentCloudTeoBindZoneToPlanCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_bind_zone_to_plan.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewBindZoneToPlanRequest()
	)

	zoneId := d.Get("zone_id").(string)
	planId := d.Get("plan_id").(string)
	request.ZoneId = helper.String(zoneId)
	request.PlanId = helper.String(planId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().BindZoneToPlanWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s bind zone to plan failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(helper.BuildToken())
	return resourceTencentCloudTeoBindZoneToPlanRead(d, meta)
}

func resourceTencentCloudTeoBindZoneToPlanRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_bind_zone_to_plan.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTeoBindZoneToPlanDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_bind_zone_to_plan.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
