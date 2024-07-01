package cvm

import (
	"context"

	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCvmImageQuotaReadOutputContent(ctx context.Context) interface{} {
	d := tccommon.ResourceDataFromContext(ctx)

	return map[string]interface{}{
		"image_num_quota": d.Get("image_num_quota"),
	}
}

func dataSourceTencentCloudCvmImageQuotaReadPostHandleResponse0(ctx context.Context, req map[string]interface{}, resp *cvm.DescribeImageQuotaResponseParams) error {
	d := tccommon.ResourceDataFromContext(ctx)

	d.SetId(helper.BuildToken())

	return nil
}
