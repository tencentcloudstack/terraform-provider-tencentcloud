package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcVpnGatewayCcnRoutesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcVpnGatewayCcnRoutes,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_vpn_gateway_ccn_routes.vpn_gateway_ccn_routes", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_vpn_gateway_ccn_routes.vpn_gateway_ccn_routes",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcVpnGatewayCcnRoutes = `

resource "tencentcloud_vpc_vpn_gateway_ccn_routes" "vpn_gateway_ccn_routes" {
  vpn_gateway_id = "vpngw-c6orbuv7"
  routes {
		route_id = "vpnr-7t3tknmg"
		status = "ENABLE"
		destination_cidr_block = "10.2.2.0/24"

  }
}

`
