package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcVpnGatewayRenewResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcVpnGatewayRenew,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_vpn_gateway_renew.vpn_gateway_renew", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_vpn_gateway_renew.vpn_gateway_renew",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcVpnGatewayRenew = `

resource "tencentcloud_vpc_vpn_gateway_renew" "vpn_gateway_renew" {
  vpn_gateway_id = "vpngw-c6orbuv7"
  instance_charge_prepaid {
		period = 1
		renew_flag = "NOTIFY_AND_AUTO_RENEW"

  }
}

`
