package teo

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
)

func resourceTencentCloudTeoOwnershipVerifyCreatePostHandleResponse0(ctx context.Context, resp *teo.VerifyOwnershipResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)

	if resp != nil && resp.Response != nil {
		_ = d.Set("status", *resp.Response.Status)
		_ = d.Set("result", *resp.Response.Result)
	}

	return nil
}
