package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixVpnGatewayRoutesDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudVpnGatewayRoutesDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpn_gateway_routes.routes"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpn_gateway_routes.routes", "vpn_gateway_route_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpn_gateway_routes.routes", "vpn_gateway_route_list.0.destination_cidr_block", "10.0.0.0/16"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpn_gateway_routes.routes", "vpn_gateway_route_list.0.instance_type", "VPNCONN"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpn_gateway_routes.routes", "vpn_gateway_route_list.0.priority", "100"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpn_gateway_routes.routes", "vpn_gateway_route_list.0.status", "ENABLE"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpn_gateway_routes.routes", "vpn_gateway_route_list.0.type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpn_gateway_routes.routes", "vpn_gateway_route_list.0.route_id"),
				),
			},
		},
	})
}

const testAccTencentCloudVpnGatewayRoutesDataSourceConfig_basic = defaultVpnDataSource + `
resource "tencentcloud_vpn_gateway_route" "route1" {
  vpn_gateway_id = data.tencentcloud_vpn_gateways.foo.gateway_list.0.id
  destination_cidr_block = "10.0.0.0/18"
  instance_type = "VPNCONN"
  instance_id = data.tencentcloud_vpn_connections.conns.connection_list.0.id
  priority = "100"
  status = "ENABLE"
}

data "tencentcloud_vpn_gateway_routes" "routes" {
  vpn_gateway_id = data.tencentcloud_vpn_gateways.foo.gateway_list.0.id
}
`
