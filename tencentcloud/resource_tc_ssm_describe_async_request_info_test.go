package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSsmDescribeAsyncRequestInfoResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSsmDescribeAsyncRequestInfo,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssm_describe_async_request_info.describe_async_request_info", "id")),
			},
			{
				ResourceName:      "tencentcloud_ssm_describe_async_request_info.describe_async_request_info",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSsmDescribeAsyncRequestInfo = `

resource "tencentcloud_ssm_describe_async_request_info" "describe_async_request_info" {
  flow_i_d = 1
}

`
