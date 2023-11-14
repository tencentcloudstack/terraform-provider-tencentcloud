package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcVpnGatewaySslClientCertResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcVpnGatewaySslClientCert,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_vpn_gateway_ssl_client_cert.vpn_gateway_ssl_client_cert", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_vpn_gateway_ssl_client_cert.vpn_gateway_ssl_client_cert",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcVpnGatewaySslClientCert = `

resource "tencentcloud_vpc_vpn_gateway_ssl_client_cert" "vpn_gateway_ssl_client_cert" {
  ssl_vpn_client_id = "vpnc-123456"
}

`
