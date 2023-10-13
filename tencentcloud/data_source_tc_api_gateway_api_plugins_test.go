package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudApiGatewayApiPluginsDataSource_basic -v
func TestAccTencentCloudApiGatewayApiPluginsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApiGatewayApiPluginsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_api_gateway_api_plugins.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_api_gateway_api_plugins.example", "api_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_api_gateway_api_plugins.example", "service_id"),
				),
			},
		},
	})
}

const testAccApiGatewayApiPluginsDataSource = `
data "tencentcloud_api_gateway_api_plugins" "example" {
  api_id     = "api-0cvmf4x4"
  service_id = "service-nxz6yync"
}
`
