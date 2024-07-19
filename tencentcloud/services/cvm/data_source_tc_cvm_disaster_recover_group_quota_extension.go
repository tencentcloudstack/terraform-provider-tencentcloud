package cvm

import (
	"context"

	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCvmDisasterRecoverGroupQuotaReadOutputContent(ctx context.Context) interface{} {
	d := tccommon.ResourceDataFromContext(ctx)

	return map[string]interface{}{
		"group_quota":             d.Get("group_quota"),
		"current_num":             d.Get("current_num"),
		"cvm_in_host_group_quota": d.Get("cvm_in_host_group_quota"),
		"cvm_in_sw_group_quota":   d.Get("cvm_in_sw_group_quota"),
		"cvm_in_rack_group_quota": d.Get("cvm_in_rack_group_quota"),
	}
}

func dataSourceTencentCloudCvmDisasterRecoverGroupQuotaReadPostHandleResponse0(ctx context.Context, req map[string]interface{}, resp *cvm.DescribeDisasterRecoverGroupQuotaResponseParams) error {
	d := tccommon.ResourceDataFromContext(ctx)

	d.SetId(helper.BuildToken())

	return nil
}
