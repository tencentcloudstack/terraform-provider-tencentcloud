package cvm

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCvmDisasterRecoverGroupQuotaReadOutputContent(ctx context.Context) interface{} {
	d := tccommon.ResourceDataFromContext(ctx)

	d.SetId(helper.BuildToken())

	return map[string]interface{}{
		"group_quota":             d.Get("group_quota"),
		"current_num":             d.Get("current_num"),
		"cvm_in_host_group_quota": d.Get("cvm_in_host_group_quota"),
		"cvm_in_sw_group_quota":   d.Get("cvm_in_sw_group_quota"),
		"cvm_in_rack_group_quota": d.Get("cvm_in_rack_group_quota"),
	}
}
