package teo

import (
	"context"

	v20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTeoSecurityIpGroupCreatePostFillRequest0(ctx context.Context, req *v20220901.CreateSecurityIPGroupRequest) error {
	req.IPGroup.GroupId = helper.IntInt64(0)
	return nil
}

func resourceTencentCloudTeoSecurityIpGroupUpdatePostFillRequest0(ctx context.Context, req *v20220901.ModifySecurityIPGroupRequest) error {
	req.Mode = helper.String("update")
	return nil
}
