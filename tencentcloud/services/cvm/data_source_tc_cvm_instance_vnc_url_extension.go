package cvm

import (
	"context"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func dataSourceTencentCloudCvmInstanceVncUrlReadOutputContent(ctx context.Context) interface{} {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}
	return map[string]interface{}{
		"instance_vnc_url": d.Get("instance_vnc_url"),
	}
}
