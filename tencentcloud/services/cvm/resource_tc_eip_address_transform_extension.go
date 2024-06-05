package cvm

import (
	"context"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"
	"time"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func resourceTencentCloudEipAddressTransformCreatePostHandleResponse0(ctx context.Context, resp *vpc.TransformAddressResponse) error {
	meta := tccommon.ProviderMetaFromContext(ctx)
	service := svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	taskId := *resp.Response.TaskId
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 1*tccommon.ReadRetryTimeout, time.Second, service.VpcIpv6AddressStateRefreshFunc(helper.UInt64ToStr(taskId), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
