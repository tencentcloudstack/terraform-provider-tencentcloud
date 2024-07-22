package teo

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func resourceTencentCloudTeoRealtimeLogDeliveryCreatePostHandleResponse0(ctx context.Context, resp *teo.CreateRealtimeLogDeliveryTaskResponse) error {
	taskId := *resp.Response.TaskId
	return checkRealtimeLogDeliveryStatus(ctx, taskId, "enabled")
}

func resourceTencentCloudTeoRealtimeLogDeliveryUpdatePostHandleResponse0(ctx context.Context, resp *teo.ModifyRealtimeLogDeliveryTaskResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}
	var taskId string
	if v, ok := d.GetOk("task_id"); ok {
		taskId = v.(string)
	}
	return checkRealtimeLogDeliveryStatus(ctx, taskId, "update")
}

func resourceTencentCloudTeoRealtimeLogDeliveryDeletePostHandleResponse0(ctx context.Context, resp *teo.DeleteRealtimeLogDeliveryTaskResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}
	var taskId string
	if v, ok := d.GetOk("task_id"); ok {
		taskId = v.(string)
	}
	return checkRealtimeLogDeliveryStatus(ctx, taskId, "deleted")
}

func checkRealtimeLogDeliveryStatus(ctx context.Context, taskId string, expectedStatuses ...string) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}
	meta := tccommon.ProviderMetaFromContext(ctx)
	if meta == nil {
		return fmt.Errorf("provider meta can not be nil")
	}

	var zoneId string
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("delivery_status"); ok && expectedStatuses[0] == "update" {
		expectedStatuses = append(expectedStatuses, v.(string))
	}

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	return resource.Retry(6*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeTeoRealtimeLogDeliveryById(ctx, zoneId, taskId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}

		if instance == nil {
			if expectedStatuses[0] != "deleted" {
				return resource.NonRetryableError(fmt.Errorf("RealtimeLogDeliveryTask data not found, taskId: %v", taskId))
			}
			return nil
		}

		for _, s := range expectedStatuses {
			if s == *instance.DeliveryStatus {
				return nil
			}
		}

		return resource.RetryableError(fmt.Errorf("RealtimeLogDeliveryTask status is %v, retry...", *instance.DeliveryStatus))
	})
}
