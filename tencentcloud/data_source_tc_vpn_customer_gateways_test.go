package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpnCustomerGatewaysDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudVpnCustomerGatewaysDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpn_customer_gateways.cgws"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpn_customer_gateways.cgws", "gateway_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpn_customer_gateways.cgws", "gateway_list.0.name", "terraform_test"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpn_customer_gateways.cgws", "gateway_list.0.public_ip_address", "1.1.1.3"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpn_customer_gateways.cgws", "gateway_list.0.tags.test", "tf"),
				),
			},
		},
	})
}

const testAccTencentCloudVpnCustomerGatewaysDataSourceConfig_basic = `
resource "tencentcloud_vpn_customer_gateway" "my_cgw" {
  name              = "terraform_test"
  public_ip_address = "1.1.1.3"
  tags = {
    test = "tf"
  }
}

data "tencentcloud_vpn_customer_gateways" "cgws" {
  id = tencentcloud_vpn_customer_gateway.my_cgw.id
}
`
