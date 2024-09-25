package tke

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	v20180525 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
)

func resourceTencentCloudKubernetesLogConfigReadRequestOnSuccess0(ctx context.Context, resp *v20180525.DescribeLogConfigsResponseParams) *resource.RetryError {
	if resp != nil {
		if resp.Message != nil && *resp.Message != "" {
			e := fmt.Errorf(*resp.Message)
			return resource.NonRetryableError(e)
		}
	}

	return nil
}

func resourceTencentCloudKubernetesLogConfigDeletePostRequest0(ctx context.Context, req *v20180525.DeleteLogConfigsRequest, resp *v20180525.DeleteLogConfigsResponse) *resource.RetryError {
	if resp != nil && resp.Response != nil {
		message := resp.Response.Message
		if message != nil && *message != "" {
			e := fmt.Errorf(*message)
			return resource.NonRetryableError(e)
		}
	}

	return nil
}
