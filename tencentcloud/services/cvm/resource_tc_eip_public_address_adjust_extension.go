package cvm

import (
	"context"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func resourceTencentCloudEipPublicAddressAdjustCreatePostHandleResponse0(ctx context.Context, resp *vpc.AdjustPublicAddressResponse) error {
	meta := tccommon.ProviderMetaFromContext(ctx)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	vpcService := svcvpc.NewVpcService(client)
	taskId := *resp.Response.TaskId
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 1*tccommon.ReadRetryTimeout, time.Second, vpcService.VpcIpv6AddressStateRefreshFunc(helper.UInt64ToStr(taskId), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}
	return nil
}
