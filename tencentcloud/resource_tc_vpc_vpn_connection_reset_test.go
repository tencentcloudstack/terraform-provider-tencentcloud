package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcVpnConnectionResetResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcVpnConnectionReset,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_vpn_connection_reset.vpn_connection_reset", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_vpn_connection_reset.vpn_connection_reset",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcVpnConnectionReset = `

resource "tencentcloud_vpc_vpn_connection_reset" "vpn_connection_reset" {
  vpn_gateway_id = "vpngw-c6orbuv7"
  vpn_connection_id = "vpnx-osftvdea"
}

`
