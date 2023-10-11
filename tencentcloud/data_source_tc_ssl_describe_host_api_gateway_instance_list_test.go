package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDescribeHostApiGatewayInstanceListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccStepSetRegion(t, "ap-nanjing")
			testAccPreCheckCommon(t, ACCOUNT_TYPE_SSL)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeHostApiGatewayInstanceListDataSource,
				Check: resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_host_api_gateway_instance_list.describe_host_api_gateway_instance_list"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_describe_host_api_gateway_instance_list.describe_host_api_gateway_instance_list", "certificate_id", "9Bpk7XOu"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_describe_host_api_gateway_instance_list.describe_host_api_gateway_instance_list", "resource_type", "apiGateway"),
				),
			},
		},
	})
}

const testAccSslDescribeHostApiGatewayInstanceListDataSource = `

data "tencentcloud_ssl_describe_host_api_gateway_instance_list" "describe_host_api_gateway_instance_list" {
  certificate_id = "9Bpk7XOu"
  resource_type = "apiGateway"
}
`
