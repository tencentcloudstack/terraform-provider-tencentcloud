package cvm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func resourceTencentCloudCvmImageSharePermissionReadPreHandleResponse0(ctx context.Context, resp *cvm.DescribeImageSharePermissionResponseParams) error {
	d := tccommon.ResourceDataFromContext(ctx)

	if len(resp.SharePermissionSet) < 1 {
		return fmt.Errorf("resource `tencentcloud_cvm_image_share_permission` read failed")
	}
	sharePermissionSet := resp.SharePermissionSet

	accountIds := make([]string, 0)
	for _, sharePermission := range sharePermissionSet {
		accountIds = append(accountIds, *sharePermission.AccountId)
	}
	_ = d.Set("account_ids", accountIds)

	return nil
}

func resourceTencentCloudCvmImageSharePermissionUpdateOnStart(ctx context.Context) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	if d.HasChange("account_ids") {
		old, new := d.GetChange("account_ids")
		oldSet := old.(*schema.Set)
		newSet := new.(*schema.Set)
		add := newSet.Difference(oldSet).List()
		remove := oldSet.Difference(newSet).List()
		if len(add) > 0 {
			addError := service.ModifyImageSharePermission(ctx, d.Id(), IMAGE_SHARE_PERMISSION_SHARE, helper.InterfacesStrings(add))
			if addError != nil {
				return addError
			}
		}
		if len(remove) > 0 {
			removeError := service.ModifyImageSharePermission(ctx, d.Id(), IMAGE_SHARE_PERMISSION_CANCEL, helper.InterfacesStrings(remove))
			if removeError != nil {
				return removeError
			}
		}
	}

	return nil
}

func resourceTencentCloudCvmImageSharePermissionDeletePostFillRequest0(ctx context.Context, req *cvm.ModifyImageSharePermissionRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var sharePermissionSet []*cvm.SharePermission

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCvmImageSharePermissionByFilter(ctx, map[string]interface{}{"ImageId": helper.String(d.Id())})
		if e != nil {
			return tccommon.RetryError(e)
		}
		sharePermissionSet = result
		return nil
	})
	if err != nil {
		return err
	}

	accountIds := make([]string, 0)
	for _, sharePermission := range sharePermissionSet {
		accountIds = append(accountIds, *sharePermission.AccountId)
	}
	req.AccountIds = helper.StringsStringsPoint(accountIds)

	return nil
}
