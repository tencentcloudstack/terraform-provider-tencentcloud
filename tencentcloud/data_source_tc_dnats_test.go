package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDataSourceDnatsBase(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSourceDnatsBase,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dnats.multi_dnats"),
					resource.TestCheckResourceAttr("data.tencentcloud_dnats.multi_dnats", "dnat_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_dnats.multi_dnats", "dnat_list.0.description", defaultInsName),
				),
			},
		},
	})
}

const testAccTencentCloudDataSourceDnatsBase = instanceCommonTestCase + `
# Create EIP 
resource "tencentcloud_eip" "eip_dev_dnat" {
  name = "${var.instance_name}"
}

resource "tencentcloud_eip" "eip_test_dnat" {
  name = "${var.instance_name}"
}

# Create NAT Gateway
resource "tencentcloud_nat_gateway" "my_nat" {
  vpc_id         = "${var.vpc_id}"
  name           = "${var.instance_name}"
  max_concurrent = 3000000
  bandwidth      = 500

  assigned_eip_set = [
    "${tencentcloud_eip.eip_dev_dnat.public_ip}",
    "${tencentcloud_eip.eip_test_dnat.public_ip}",
  ]
}

# Add DNAT Entry
resource "tencentcloud_dnat" "dev_dnat" {
  vpc_id       = "${tencentcloud_nat_gateway.my_nat.vpc_id}"
  nat_id       = "${tencentcloud_nat_gateway.my_nat.id}"
  protocol     = "TCP"
  elastic_ip   = "${tencentcloud_eip.eip_dev_dnat.public_ip}"
  elastic_port = "80"
  private_ip   = "${tencentcloud_instance.default.private_ip}"
  private_port = "9001"
  description  = "${var.instance_name}"
}

data "tencentcloud_dnats" "multi_dnats" {
  nat_id = "${tencentcloud_dnat.dev_dnat.nat_id}"
  vpc_id = "${tencentcloud_dnat.dev_dnat.vpc_id}"
}
`
