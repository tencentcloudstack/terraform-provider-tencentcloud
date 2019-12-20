package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNatGatewaysDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudNatGatewaysDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_nat_gateways.multi_nat"),
					resource.TestCheckResourceAttr("data.tencentcloud_nat_gateways.multi_nat", "nats.#", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_nat_gateways.multi_nat", "nats.0.name", "terraform_test_nats"),
					resource.TestCheckResourceAttr("data.tencentcloud_nat_gateways.multi_nat", "nats.1.bandwidth", "500"),
				),
			},
		},
	})
}

const testAccTencentCloudNatGatewaysDataSourceConfig_basic = `
resource "tencentcloud_vpc" "main" {
  name       = "terraform_test_nats"
  cidr_block = "10.6.0.0/16"
}
resource "tencentcloud_eip" "eip_dev_dnat" {
  name = "terraform_test"
}
resource "tencentcloud_eip" "eip_test_dnat" {
  name = "terraform_test"
}

resource "tencentcloud_nat_gateway" "dev_nat" {
  vpc_id           = "${tencentcloud_vpc.main.id}"
  name             = "terraform_test_nats"
  max_concurrent   = 3000000
  bandwidth        = 500
  assigned_eip_set = [
    "${tencentcloud_eip.eip_dev_dnat.public_ip}",
  ]
}
resource "tencentcloud_nat_gateway" "test_nat" {
  vpc_id           = "${tencentcloud_vpc.main.id}"
  name             = "terraform_test_nats"
  max_concurrent   = 3000000
  bandwidth        = 500
  assigned_eip_set = [
    "${tencentcloud_eip.eip_test_dnat.public_ip}",
  ]
}

data "tencentcloud_nat_gateways" "multi_nat" {
  name           = "${tencentcloud_nat_gateway.dev_nat.name}"
  vpc_id         = "${tencentcloud_vpc.main.id}"
}
`
