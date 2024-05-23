package teo

import (
	"context"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
)

func resourceTencentCloudTeoZoneCreateRequestOnError0(ctx context.Context, req *teo.CreateZoneRequest, e error) *resource.RetryError {
	if tccommon.IsExpectError(e, []string{"ResourceInUse", "ResourceInUse.Others"}) {
		return resource.NonRetryableError(e)
	}
	return nil
}

func resourceTencentCloudTeoZoneCreatePostHandleResponse0(ctx context.Context, resp *teo.CreateZoneResponse) error {
	meta := tccommon.ProviderMetaFromContext(ctx)
	d := tccommon.ResourceDataFromContext(ctx)

	zoneId := *resp.Response.ZoneId

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err := resource.Retry(6*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeTeoZone(ctx, zoneId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if *instance.Status == "pending" {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("zone status is %v, retry...", *instance.Status))
	})
	if err != nil {
		return err
	}

	if v, _ := d.GetOkExists("paused"); v != nil {
		if v.(bool) {
			err := service.ModifyZoneStatus(ctx, zoneId, v.(bool), "create")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func resourceTencentCloudTeoZoneReadPostHandleResponse0(ctx context.Context, resp *teo.Zone) error {
	d := tccommon.ResourceDataFromContext(ctx)

	if resp.Resources != nil && len(resp.Resources) > 0 {
		if resp.Resources[0].PlanId != nil {
			_ = d.Set("plan_id", resp.Resources[0].PlanId)
		}
	}
	return nil
}

func resourceTencentCloudTeoZoneUpdatePostRequest1(ctx context.Context, req *teo.ModifyZoneStatusRequest, resp *teo.ModifyZoneStatusResponse) *resource.RetryError {
	meta := tccommon.ProviderMetaFromContext(ctx)
	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	zoneId := *req.ZoneId
	operate := "update"

	err := resource.Retry(6*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeTeoZone(ctx, zoneId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if operate == "delete" {
			if *instance.ActiveStatus == "paused" {
				return nil
			}
		} else {
			if *instance.ActiveStatus == "inactive" || *instance.ActiveStatus == "paused" {
				return nil
			}
		}
		return resource.RetryableError(fmt.Errorf("zone status is %v, retry...", *instance.ActiveStatus))
	})
	if err != nil {
		return resource.NonRetryableError(err)
	}

	return nil
}

func resourceTencentCloudTeoZoneDeletePostFillRequest0(ctx context.Context, req *teo.DeleteZoneRequest) error {
	meta := tccommon.ProviderMetaFromContext(ctx)
	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	zoneId := *req.ZoneId

	instance, err := service.DescribeTeoZone(ctx, zoneId)
	if err != nil {
		return err
	}

	if !*instance.Paused {
		err := service.ModifyZoneStatus(ctx, zoneId, true, "delete")
		if err != nil {
			return err
		}
	}
	return nil
}
