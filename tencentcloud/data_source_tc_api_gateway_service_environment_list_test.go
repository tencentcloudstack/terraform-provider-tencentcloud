package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudApiGatewayServiceEnvironmentListDataSource_basic -v
func TestAccTencentCloudApiGatewayServiceEnvironmentListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApiGatewayServiceEnvironmentListDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_api_gateway_service_environment_list.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_api_gateway_service_environment_list.example", "service_id"),
				),
			},
		},
	})
}

const testAccApiGatewayServiceEnvironmentListDataSource = `
data "tencentcloud_api_gateway_service_environment_list" "example" {
  service_id = "service-nxz6yync"
}
`
