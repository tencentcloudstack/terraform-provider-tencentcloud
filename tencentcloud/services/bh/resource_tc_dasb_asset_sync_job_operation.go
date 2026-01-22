package bh

import (
	"fmt"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dasb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dasb/v20191018"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDasbAssetSyncJobOperationOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDasbAssetSyncJobOperationCreate,
		Read:   resourceTencentCloudDasbAssetSyncJobOperationRead,
		Delete: resourceTencentCloudDasbAssetSyncJobOperationDelete,

		Schema: map[string]*schema.Schema{
			"category": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Synchronize asset categories, 1- Host assets, 2- Database assets.",
			},
		},
	}
}

func resourceTencentCloudDasbAssetSyncJobOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_asset_sync_job_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		request  = dasb.NewCreateAssetSyncJobRequest()
		waitReq  = dasb.NewDescribeAssetSyncStatusRequest()
		category string
	)

	if v, ok := d.GetOkExists("category"); ok {
		request.Category = helper.IntUint64(v.(int))
		waitReq.Category = helper.IntUint64(v.(int))
		category = strconv.Itoa(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDasbClient().CreateAssetSyncJob(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dasb AssetSyncJob failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(category)

	// wait
	err = resource.Retry(4*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDasbClient().DescribeAssetSyncStatus(waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe dasb AssetSyncJob failed, Response is nil."))
		}

		if result.Response.Status.InProcess == nil {
			return resource.NonRetryableError(fmt.Errorf("InProcess is nil."))
		}

		if !*result.Response.Status.InProcess {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Dasb asset sync job is still running..."))
	})

	if err != nil {
		log.Printf("[CRITAL]%s describe dasb AssetSyncJob failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDasbAssetSyncJobOperationRead(d, meta)
}

func resourceTencentCloudDasbAssetSyncJobOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_asset_sync_job_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDasbAssetSyncJobOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_asset_sync_job_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
