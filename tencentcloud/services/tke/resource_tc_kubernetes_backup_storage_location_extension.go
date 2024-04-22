package tke

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func resourceTencentCloudKubernetesBackupStorageLocationCreatePostHandleResponse0(ctx context.Context, resp *tke.CreateBackupStorageLocationResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}

	meta := tccommon.ProviderMetaFromContext(ctx)
	if meta == nil {
		return fmt.Errorf("provider meta can not be nil")
	}
	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	// wait for status ok
	return resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		locations, errRet := service.DescribeBackupStorageLocations(ctx, []string{d.Get("name").(string)})
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if len(locations) != 1 {
			resource.RetryableError(fmt.Errorf("more than 1 location returnen in api response, expected 1 but got %d", len(locations)))
		}
		if locations[0].State == nil {
			return resource.RetryableError(fmt.Errorf("location %s is still in state nil", d.Get("name").(string)))
		}
		if len(locations) == 1 && *locations[0].State == backupStorageLocationStateAvailable {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("location %s is still in state %s", d.Get("name").(string), *locations[0].State))
	})
}

func resourceTencentCloudKubernetesBackupStorageLocationDeletePostHandleResponse0(ctx context.Context, resp *tke.DeleteBackupStorageLocationResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}
	meta := tccommon.ProviderMetaFromContext(ctx)
	if meta == nil {
		return fmt.Errorf("provider meta can not be nil")
	}
	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	// wait until location is deleted
	return resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		locations, errRet := service.DescribeBackupStorageLocations(ctx, []string{d.Id()})
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if len(locations) == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("location %s is still not deleted", d.Id()))
	})
}
