package tke

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tkev20180525 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudKubernetesCancelUpgradePlanOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKubernetesCancelUpgradePlanOperationCreate,
		Read:   resourceTencentCloudKubernetesCancelUpgradePlanOperationRead,
		Delete: resourceTencentCloudKubernetesCancelUpgradePlanOperationDelete,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster ID.",
			},

			"plan_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Upgrade plan ID.",
			},
		},
	}
}

func resourceTencentCloudKubernetesCancelUpgradePlanOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cancel_upgrade_plan_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = tkev20180525.NewCancelUpgradePlanRequest()
		clusterId string
		planId    string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterID = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOkExists("plan_id"); ok {
		request.PlanID = helper.IntInt64(v.(int))
		planId = helper.IntToStr(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().CancelUpgradePlanWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create kubernetes cancel upgrade plan operation failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(strings.Join([]string{clusterId, planId}, tccommon.FILED_SP))
	return resourceTencentCloudKubernetesCancelUpgradePlanOperationRead(d, meta)
}

func resourceTencentCloudKubernetesCancelUpgradePlanOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cancel_upgrade_plan_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudKubernetesCancelUpgradePlanOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cancel_upgrade_plan_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
