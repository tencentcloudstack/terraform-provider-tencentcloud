package tse_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTseGatewayRoutesDataSource_basic -v
func TestAccTencentCloudTseGatewayRoutesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseGatewayRoutesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tse_gateway_routes.gateway_routes"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.route_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.route_list.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.route_list.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.route_list.0.methods.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.route_list.0.paths.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.route_list.0.hosts.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.route_list.0.protocols.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.route_list.0.preserve_host"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.route_list.0.https_redirect_status_code"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.route_list.0.strip_path"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.route_list.0.created_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.route_list.0.destination_ports.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.route_list.0.service_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.route_list.0.service_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.route_list.0.headers.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.route_list.0.headers.0.key"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_routes.gateway_routes", "result.0.route_list.0.headers.0.value"),
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
