package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceTencentCloudVpcV3Subnets_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudVpcRouteTables,

				Check: resource.ComposeTestCheckFunc(
					//id filter
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_route_tables.id_instances"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpc_route_tables.id_instances", "instance_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpc_route_tables.id_instances", "instance_list.0.name", "ci-temp-test-rt"),

					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.id_instances", "instance_list.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.id_instances", "instance_list.0.route_table_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.id_instances", "instance_list.0.is_default"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.id_instances", "instance_list.0.subnet_ids.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.id_instances", "instance_list.0.route_entry_infos.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.id_instances", "instance_list.0.create_time"),

					//name filter ,Every routable with a "ci-temp-test-rt" name will be found
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_route_tables.name_instances"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.name_instances", "instance_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.id_instances", "instance_list.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.id_instances", "instance_list.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.id_instances", "instance_list.0.route_table_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.id_instances", "instance_list.0.is_default"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.id_instances", "instance_list.0.subnet_ids.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.id_instances", "instance_list.0.route_entry_infos.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.id_instances", "instance_list.0.create_time"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudVpcRouteTables = `
variable "availability_zone" {
	default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
    name="guagua-ci-temp-test"
    cidr_block="10.0.0.0/16"
}

resource "tencentcloud_route_table" "route_table" {
   vpc_id = "${tencentcloud_vpc.foo.id}"
   name = "ci-temp-test-rt"
}

data "tencentcloud_vpc_route_tables" "id_instances" {
	route_table_id="${tencentcloud_route_table.route_table.id}"
}
data "tencentcloud_vpc_route_tables" "name_instances" {
	name="${tencentcloud_route_table.route_table.name}"
}

`
