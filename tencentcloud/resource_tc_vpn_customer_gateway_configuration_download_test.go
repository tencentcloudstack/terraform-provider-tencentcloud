package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpnCustomerGatewayConfigurationDownloadResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcVpnCustomerGatewayConfigurationDownload,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpn_customer_gateway_configuration_download.vpn_customer_gateway_configuration_download", "id")),
			},
		},
	})
}

const testAccVpcVpnCustomerGatewayConfigurationDownload = `

resource "tencentcloud_vpn_customer_gateway_configuration_download" "vpn_customer_gateway_configuration_download" {
  vpn_gateway_id    = "vpngw-gt8bianl"
  vpn_connection_id = "vpnx-kme2tx8m"
  customer_gateway_vendor {
    platform         = "comware"
    software_version = "V1.0"
    vendor_name      = "h3c"
  }
  interface_name    = "test"
}

`
