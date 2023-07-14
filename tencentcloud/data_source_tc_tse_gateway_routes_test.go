package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixTseGatewayRoutesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseGatewayRoutesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tse_gateway_routes.gateway_routes"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.RouteList.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.RouteList.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.RouteList.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.RouteList.0.methods.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.RouteList.0.paths.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.RouteList.0.hosts.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.RouteList.0.protocols.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.RouteList.0.preserveHost"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.RouteList.0.https_redirect_status_code"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.RouteList.0.strip_path"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.RouteList.0.created_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.RouteList.0.destination_ports"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.RouteList.0.service_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.RouteList.0.service_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.RouteList.0.headers.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.RouteList.0.headers.0.key"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.RouteList.0.headers.0.value"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.total_count"),
				),
			},
		},
	})
}

const testAccTseGatewayRoutesDataSource = `

data "tencentcloud_tse_gateway_routes" "gateway_routes" {
	gateway_id   = "gateway-ddbb709b"
	service_name = "test"
	route_name   = "keep-routes"
}

`
