package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTencentCloudVpcV3RouteTables_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudVpcRouteTables,
				Check: resource.ComposeTestCheckFunc(
					// id filter
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_route_tables.id_instances"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpc_route_tables.id_instances", "instance_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpc_route_tables.id_instances", "instance_list.0.name", "ci-temp-test-rt"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.id_instances", "instance_list.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.id_instances", "instance_list.0.route_table_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.id_instances", "instance_list.0.is_default"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.id_instances", "instance_list.0.subnet_ids.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.id_instances", "instance_list.0.route_entry_infos.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.id_instances", "instance_list.0.create_time"),

					// name filter ,Every routable with a "ci-temp-test-rt" name will be found
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_route_tables.name_instances"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.name_instances", "instance_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.name_instances", "instance_list.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.name_instances", "instance_list.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.name_instances", "instance_list.0.route_table_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.name_instances", "instance_list.0.is_default"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.name_instances", "instance_list.0.subnet_ids.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.name_instances", "instance_list.0.route_entry_infos.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.name_instances", "instance_list.0.create_time"),

					// tags filter ,Every routable with a tag test:test will be found
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_route_tables.tags_instances"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.tags_instances", "instance_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.tags_instances", "instance_list.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.tags_instances", "instance_list.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.tags_instances", "instance_list.0.route_table_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.tags_instances", "instance_list.0.is_default"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.tags_instances", "instance_list.0.subnet_ids.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.tags_instances", "instance_list.0.route_entry_infos.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.tags_instances", "instance_list.0.create_time"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpc_route_tables.tags_instances", "instance_list.0.tags.test", "test"),

					// vpc_id filter
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_route_tables.vpc_instances"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.vpc_instances", "instance_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.tags_instances", "instance_list.0.vpc_id"),

					// vpc_id && association_main filter
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_route_tables.vpc_default_instance"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpc_route_tables.vpc_default_instance", "instance_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vpc_route_tables.vpc_default_instance", "instance_list.0.vpc_id"),
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
  name       = "guagua-ci-temp-test"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_route_table" "route_table" {
  vpc_id = tencentcloud_vpc.foo.id
  name   = "ci-temp-test-rt"

  tags = {
    "test" = "test"
  }
}

data "tencentcloud_vpc_route_tables" "id_instances" {
  route_table_id = tencentcloud_route_table.route_table.id
}

data "tencentcloud_vpc_route_tables" "name_instances" {
  name = tencentcloud_route_table.route_table.name
}

data "tencentcloud_vpc_route_tables" "vpc_instances" {
  vpc_id = tencentcloud_vpc.foo.id
}

data "tencentcloud_vpc_route_tables" "vpc_default_instance" {
  vpc_id           = tencentcloud_vpc.foo.id
  association_main = true
}

data "tencentcloud_vpc_route_tables" "tags_instances" {
  tags = tencentcloud_route_table.route_table.tags
}
`
