package teo

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func resourceTencentCloudTeoAccelerationDomainCreatePostHandleResponse0(ctx context.Context, resp *teo.CreateAccelerationDomainResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}
	return checkAccelerationDomainStatus(ctx, d.Timeout(schema.TimeoutCreate), "online")
}

func resourceTencentCloudTeoAccelerationDomainUpdateOnExit(ctx context.Context) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}
	status := "online"
	if v, ok := d.GetOk("status"); ok {
		status = v.(string)
	}
	return checkAccelerationDomainStatus(ctx, d.Timeout(schema.TimeoutUpdate), status)
}

func resourceTencentCloudTeoAccelerationDomainDeletePostHandleResponse0(ctx context.Context, resp *teo.ModifyAccelerationDomainStatusesResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}
	return checkAccelerationDomainStatus(ctx, d.Timeout(schema.TimeoutDelete), "offline")
}

func checkAccelerationDomainStatus(ctx context.Context, retryTimeout time.Duration, expectedStatuses ...string) error {
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
	var domainName string
	if v, ok := d.GetOk("domain_name"); ok {
		domainName = v.(string)
	}

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	return resource.Retry(retryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeTeoAccelerationDomainById(ctx, zoneId, domainName)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}

		for _, s := range expectedStatuses {
			if s == *instance.DomainStatus {
				return nil
			}
		}

		return resource.RetryableError(fmt.Errorf("AccelerationDomain status is %v, retry...", *instance.DomainStatus))
	})
}
