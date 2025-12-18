package bh

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bhv20230418 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bh/v20230418"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudBhUserSyncTaskOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBhUserSyncTaskOperationCreate,
		Read:   resourceTencentCloudBhUserSyncTaskOperationRead,
		Delete: resourceTencentCloudBhUserSyncTaskOperationDelete,
		Schema: map[string]*schema.Schema{
			"user_kind": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Synchronized user type, 1-synchronize IOA users.",
			},
		},
	}
}

func resourceTencentCloudBhUserSyncTaskOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_user_sync_task_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = bhv20230418.NewCreateSyncUserTaskRequest()
		userKind string
	)

	if v, ok := d.GetOkExists("user_kind"); ok {
		request.UserKind = helper.IntUint64(v.(int))
		userKind = helper.IntToStr(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().CreateSyncUserTaskWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create bh user sync task failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(userKind)

	// wait
	waitReq := bhv20230418.NewDescribeUserSyncStatusRequest()
	if v, ok := d.GetOkExists("user_kind"); ok {
		waitReq.UserKind = helper.IntUint64(v.(int))
	}

	reqErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().DescribeUserSyncStatusWithContext(ctx, waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s describe bh user sync task failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudBhUserSyncTaskOperationRead(d, meta)
}

func resourceTencentCloudBhUserSyncTaskOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_user_sync_task_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudBhUserSyncTaskOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_user_sync_task_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
