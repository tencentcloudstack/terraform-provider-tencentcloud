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

func ResourceTencentCloudBhSyncDevicesToIoaOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBhSyncDevicesToIoaOperationCreate,
		Read:   resourceTencentCloudBhSyncDevicesToIoaOperationRead,
		Delete: resourceTencentCloudBhSyncDevicesToIoaOperationDelete,
		Schema: map[string]*schema.Schema{
			"device_id_set": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Description: "Asset ID collection. Assets must be bound to bastion host instances that support IOA functionality. Maximum 200 assets can be synchronized at a time.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func resourceTencentCloudBhSyncDevicesToIoaOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_sync_devices_to_ioa_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = bhv20230418.NewSyncDevicesToIOARequest()
		deviceIds []string
	)

	if v, ok := d.GetOk("device_id_set"); ok {
		deviceIdSetSet := v.(*schema.Set).List()
		for i := range deviceIdSetSet {
			deviceIdSet := deviceIdSetSet[i].(int)
			request.DeviceIdSet = append(request.DeviceIdSet, helper.IntUint64(deviceIdSet))
			deviceIds = append(deviceIds, helper.IntToStr(deviceIdSet))
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().SyncDevicesToIOAWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create bh sync devices to ioa operation failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(helper.HashStrings(deviceIds))
	return resourceTencentCloudBhSyncDevicesToIoaOperationRead(d, meta)
}

func resourceTencentCloudBhSyncDevicesToIoaOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_sync_devices_to_ioa_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudBhSyncDevicesToIoaOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_sync_devices_to_ioa_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
