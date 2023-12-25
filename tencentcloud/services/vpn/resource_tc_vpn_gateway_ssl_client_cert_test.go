package vpn_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpnGatewaySslClientCertResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnGatewaySslClientCert,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpn_gateway_ssl_client_cert.vpn_gateway_ssl_client_cert", "id")),
			},
			{
				Config: testAccVpnGatewaySslClientCertUpdate,
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

const testAccVpnGatewaySslClientCert = `

resource "tencentcloud_vpn_ssl_server" "server" {
  local_address       = [
    "172.16.0.0/17",
  ]
  remote_address      = "173.16.1.0/24"
  ssl_vpn_server_name = "tf-vpn-ssl-server"
  vpn_gateway_id      = "vpngw-mll9np1x"
  ssl_vpn_protocol = "UDP"
  ssl_vpn_port = 1194
  integrity_algorithm = "MD5"
  encrypt_algorithm = "AES-128-CBC"
  compress = false
}

resource "tencentcloud_vpn_ssl_client" "client" {
  ssl_vpn_server_id = tencentcloud_vpn_ssl_server.server.id
  ssl_vpn_client_name = "tf-vpn-ssl-client"
}

resource "tencentcloud_vpn_gateway_ssl_client_cert" "vpn_gateway_ssl_client_cert" {
  ssl_vpn_client_id = tencentcloud_vpn_ssl_client.client.id
  switch = "off"
}

`

const testAccVpnGatewaySslClientCertUpdate = `

resource "tencentcloud_vpn_ssl_server" "server" {
  local_address       = [
    "172.16.0.0/17",
  ]
  remote_address      = "173.16.1.0/24"
  ssl_vpn_server_name = "tf-vpn-ssl-server"
  vpn_gateway_id      = "vpngw-mll9np1x"
  ssl_vpn_protocol = "UDP"
  ssl_vpn_port = 1194
  integrity_algorithm = "MD5"
  encrypt_algorithm = "AES-128-CBC"
  compress = false
}

resource "tencentcloud_vpn_ssl_client" "client" {
  ssl_vpn_server_id = tencentcloud_vpn_ssl_server.server.id
  ssl_vpn_client_name = "tf-vpn-ssl-client"
}

resource "tencentcloud_vpn_gateway_ssl_client_cert" "vpn_gateway_ssl_client_cert" {
  ssl_vpn_client_id = tencentcloud_vpn_ssl_client.client.id
  switch = "on"
}

`
