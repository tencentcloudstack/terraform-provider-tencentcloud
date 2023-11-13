package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcVpnCustomerGatewayConfigurationDownloadResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcVpnCustomerGatewayConfigurationDownload,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_vpn_customer_gateway_configuration_download.vpn_customer_gateway_configuration_download", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_vpn_customer_gateway_configuration_download.vpn_customer_gateway_configuration_download",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcVpnCustomerGatewayConfigurationDownload = `

resource "tencentcloud_vpc_vpn_customer_gateway_configuration_download" "vpn_customer_gateway_configuration_download" {
  vpn_gateway_id = "vpngw-c6orbuv7"
  vpn_connection_id = "vpnx-osftvdea"
  customer_gateway_vendor {
		platform = "comware"
		software_version = "V1.0"
		vendor_name = "h3c"

  }
  interface_name = ""
}

`
