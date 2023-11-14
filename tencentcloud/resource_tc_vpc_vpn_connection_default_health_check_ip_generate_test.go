package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcVpnConnectionDefaultHealthCheckIpGenerateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcVpnConnectionDefaultHealthCheckIpGenerate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_vpn_connection_default_health_check_ip_generate.vpn_connection_default_health_check_ip_generate", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_vpn_connection_default_health_check_ip_generate.vpn_connection_default_health_check_ip_generate",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcVpnConnectionDefaultHealthCheckIpGenerate = `

resource "tencentcloud_vpc_vpn_connection_default_health_check_ip_generate" "vpn_connection_default_health_check_ip_generate" {
  vpn_gateway_id = "vpngw-c6orbuv7"
}

`
