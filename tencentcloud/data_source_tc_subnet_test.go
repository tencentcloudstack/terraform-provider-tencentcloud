package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceTencentCloudSubnet_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudSubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_subnet.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_subnet.foo", "name", "guagua_vpc_subnet_test"),
					resource.TestCheckResourceAttr("data.tencentcloud_subnet.foo", "availability_zone", "ap-guangzhou-3"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudSubnetConfig = `
variable "availability_zone" {
	default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
    name       = "guagua-ci-temp-test"
    cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
	availability_zone = "${var.availability_zone}"
	name              = "guagua_vpc_subnet_test"
	vpc_id            = "${tencentcloud_vpc.foo.id}"
	cidr_block        = "10.0.20.0/28"
	is_multicast      = false
}

data "tencentcloud_subnet" "foo" {
	vpc_id    = "${tencentcloud_vpc.foo.id}"
	subnet_id = "${tencentcloud_subnet.subnet.id}"
}
`
