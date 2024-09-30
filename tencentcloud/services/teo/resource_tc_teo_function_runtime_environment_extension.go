package teo

import (
	"context"

	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTeoFunctionRuntimeEnvironmentCreatePostFillRequest0(ctx context.Context, req *teov20220901.HandleFunctionRuntimeEnvironmentRequest) error {
	req.Operation = helper.String("setEnvironmentVariable")
	return nil
}

func resourceTencentCloudTeoFunctionRuntimeEnvironmentUpdatePostFillRequest0(ctx context.Context, req *teov20220901.HandleFunctionRuntimeEnvironmentRequest) error {
	req.Operation = helper.String("resetEnvironmentVariable")
	return nil
}

func resourceTencentCloudTeoFunctionRuntimeEnvironmentDeletePostFillRequest0(ctx context.Context, req *teov20220901.HandleFunctionRuntimeEnvironmentRequest) error {
	req.Operation = helper.String("deleteEnvironmentVariable")
	return nil
}
