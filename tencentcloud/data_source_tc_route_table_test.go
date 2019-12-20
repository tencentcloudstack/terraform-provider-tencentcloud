package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTencentCloudRouteTable_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTencentCloudRouteTableConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_route_table.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_route_table.foo", "name", "tf-ci-test"),
				),
			},
		},
	})
}

const testAccDataSourceTencentCloudRouteTableConfig = `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
  name       = "tf-ci-test"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_route_table" "route_table" {
  vpc_id = "${tencentcloud_vpc.foo.id}"
  name   = "tf-ci-test"
}

data "tencentcloud_route_table" "foo" {
  route_table_id = "${tencentcloud_route_table.route_table.id}"
}
`
