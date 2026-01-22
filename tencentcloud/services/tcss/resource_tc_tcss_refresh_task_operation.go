package tcss

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcssv20201101 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcss/v20201101"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTcssRefreshTaskOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcssRefreshTaskOperationCreate,
		Read:   resourceTencentCloudTcssRefreshTaskOperationRead,
		Delete: resourceTencentCloudTcssRefreshTaskOperationDelete,
		Schema: map[string]*schema.Schema{
			"cluster_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Description: "Cluster Id list.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"is_sync_list_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to sync list only.",
			},
		},
	}
}

func resourceTencentCloudTcssRefreshTaskOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcss_refresh_task_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = tcssv20201101.NewCreateRefreshTaskRequest()
	)

	if v, ok := d.GetOk("cluster_ids"); ok {
		clusterIDsSet := v.(*schema.Set).List()
		for i := range clusterIDsSet {
			if clusterId, ok := clusterIDsSet[i].(string); ok {
				request.ClusterIDs = append(request.ClusterIDs, helper.String(clusterId))
			}
		}
	}

	if v, ok := d.GetOkExists("is_sync_list_only"); ok {
		request.IsSyncListOnly = helper.Bool(v.(bool))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTcssV20201101Client().CreateRefreshTaskWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create tcss refresh task operation failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(helper.BuildToken())
	return resourceTencentCloudTcssRefreshTaskOperationRead(d, meta)
}

func resourceTencentCloudTcssRefreshTaskOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcss_refresh_task_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTcssRefreshTaskOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcss_refresh_task_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
