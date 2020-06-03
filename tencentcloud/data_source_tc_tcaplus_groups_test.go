package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataTcaplusGroupsName = "data.tencentcloud_tcaplus_groups.id_test"

func TestAccTencentCloudDataTcaplusGroups(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTcaplusGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataTcaplusGroupsBaic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTcaplusGroupExists("tencentcloud_tcaplus_group.test_group"),
					resource.TestCheckResourceAttrSet(testDataTcaplusGroupsName, "cluster_id"),
					resource.TestCheckResourceAttrSet(testDataTcaplusGroupsName, "tablegroup_id"),
					resource.TestCheckResourceAttr(testDataTcaplusGroupsName, "list.#", "1"),
					resource.TestCheckResourceAttr(testDataTcaplusGroupsName, "list.0.tablegroup_name", "tf_test_group_name_guagua"),
					resource.TestCheckResourceAttr(testDataTcaplusGroupsName, "list.0.table_count", "0"),
					resource.TestCheckResourceAttrSet(testDataTcaplusGroupsName, "list.0.tablegroup_id"),
					resource.TestCheckResourceAttrSet(testDataTcaplusGroupsName, "list.0.total_size"),
					resource.TestCheckResourceAttrSet(testDataTcaplusGroupsName, "list.0.create_time"),
				),
			},
		},
	})
}

const testAccTencentCloudDataTcaplusGroupsBaic = `
variable "availability_zone" {
default = "ap-shanghai-2"
}

data "tencentcloud_vpc_subnets" "vpc" {
    is_default        = true
    availability_zone = var.availability_zone
}
resource "tencentcloud_tcaplus_cluster" "test_cluster" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_tcaplus_data_guagua"
  vpc_id                   = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_id                = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id
  password                 = "1qaA2k1wgvfa3ZZZ"
  old_password_expire_last = 3600
}
resource "tencentcloud_tcaplus_tablegroup" "test_group" {
  cluster_id       = tencentcloud_tcaplus_cluster.test_cluster.id
  tablegroup_name  = "tf_test_group_name_guagua"
}

data "tencentcloud_tcaplus_groups" "id_test" {
   cluster_id         = tencentcloud_tcaplus_cluster.test_cluster.id
   tablegroup_id      = tencentcloud_tcaplus_group.test_group.id
}
`
