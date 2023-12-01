package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpnConnectionResetResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnConnectionReset,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpn_connection_reset.vpn_connection_reset", "id")),
			},
		},
	})
}

const testAccVpnConnectionReset = `

resource "tencentcloud_vpn_connection_reset" "vpn_connection_reset" {
  vpn_gateway_id    = "vpngw-gt8bianl"
  vpn_connection_id = "vpnx-kme2tx8m"
}

`
