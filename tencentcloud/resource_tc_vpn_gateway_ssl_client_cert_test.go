package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpn_gateway_ssl_client_cert.vpn_gateway_ssl_client_cert", "id")),
			},
			{
				Config: testAccVpcVpnGatewaySslClientCertUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_gateway_ssl_client_cert.vpn_gateway_ssl_client_cert", "id"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway_ssl_client_cert.vpn_gateway_ssl_client_cert", "switch", "on"),
				),
			},
			{
				ResourceName:      "tencentcloud_vpn_gateway_ssl_client_cert.vpn_gateway_ssl_client_cert",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcVpnGatewaySslClientCert = `

resource "tencentcloud_vpn_gateway_ssl_client_cert" "vpn_gateway_ssl_client_cert" {
  ssl_vpn_client_id = "vpnc-52f5lnd5"
  switch = "off"
}

`

const testAccVpcVpnGatewaySslClientCertUpdate = `

resource "tencentcloud_vpn_gateway_ssl_client_cert" "vpn_gateway_ssl_client_cert" {
  ssl_vpn_client_id = "vpnc-52f5lnd5"
  switch = "on"
}

`
