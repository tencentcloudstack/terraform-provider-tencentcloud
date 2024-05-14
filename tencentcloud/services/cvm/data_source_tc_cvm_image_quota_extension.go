package cvm

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCvmImageQuotaReadOutputContent(ctx context.Context) interface{} {
	d := tccommon.ResourceDataFromContext(ctx)

	d.SetId(helper.BuildToken())

	return map[string]interface{}{
		"image_num_quota": d.Get("image_num_quota"),
	}
}
