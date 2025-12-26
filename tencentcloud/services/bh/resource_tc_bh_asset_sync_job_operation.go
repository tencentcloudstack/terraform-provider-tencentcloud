package bh

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bhv20230418 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bh/v20230418"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudBhAssetSyncJobOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBhAssetSyncJobOperationCreate,
		Read:   resourceTencentCloudBhAssetSyncJobOperationRead,
		Delete: resourceTencentCloudBhAssetSyncJobOperationDelete,
		Schema: map[string]*schema.Schema{
			"category": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{1, 2, 3}),
				Description:  "Asset synchronization category. 1 - host assets, 2 - database assets, 3 - Container assets.",
			},
		},
	}
}

func resourceTencentCloudBhAssetSyncJobOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_asset_sync_job_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = bhv20230418.NewCreateAssetSyncJobRequest()
		category string
	)

	if v, ok := d.GetOkExists("category"); ok {
		request.Category = helper.IntUint64(v.(int))
		category = helper.IntToStr(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().CreateAssetSyncJobWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create bh asset sync job operation failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(category)

	// wait
	waitReq := bhv20230418.NewDescribeAssetSyncStatusRequest()
	if v, ok := d.GetOkExists("category"); ok {
		waitReq.Category = helper.IntUint64(v.(int))
	}

	reqErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().DescribeAssetSyncStatusWithContext(ctx, waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Status == nil || result.Response.Status.InProcess == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe bh asset sync status failed, Response is nil."))
		}

		if !*result.Response.Status.InProcess {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Asset sync status is still running..."))
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s describe bh asset sync status failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudBhAssetSyncJobOperationRead(d, meta)
}

func resourceTencentCloudBhAssetSyncJobOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_asset_sync_job_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudBhAssetSyncJobOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_asset_sync_job_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
