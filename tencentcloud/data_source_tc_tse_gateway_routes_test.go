package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTseGatewayRoutesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseGatewayRoutesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tse_gateway_routes.gateway_routes")),
			},
		},
	})
}

const testAccTseGatewayRoutesDataSource = `

data "tencentcloud_tse_gateway_routes" "gateway_routes" {
  gateway_id = "gateway-xxxxxx"
  service_name = "serviceA"
  route_name = "123"
  filters {
		key = "name"
		value = "123"

  }
  }

`
