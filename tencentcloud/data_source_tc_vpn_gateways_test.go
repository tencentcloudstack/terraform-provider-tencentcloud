package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpnGatewaysDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudVpnGatewaysDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpn_gateways.cgws"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpn_gateways.cgws", "gateway_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpn_gateways.cgws", "gateway_list.0.name", "terraform_test"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpn_gateways.cgws", "gateway_list.0.bandwidth", "10"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpn_gateways.cgws", "gateway_list.0.tags.test", "tf"),
				),
			},
		},
	})
}

const testAccTencentCloudVpnGatewaysDataSourceConfig_basic = `
# Create VPC
data "tencentcloud_vpc_instances" "foo" {
  name = "Default-VPC"
}

resource "tencentcloud_vpn_gateway" "my_cgw" {
  name      = "terraform_test"
  vpc_id    = data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id
  bandwidth = 10
  zone      = "ap-guangzhou-3"

  tags = {
    test = "tf"
  }
}

data "tencentcloud_vpn_gateways" "cgws" {
  id = tencentcloud_vpn_gateway.my_cgw.id
}
`
