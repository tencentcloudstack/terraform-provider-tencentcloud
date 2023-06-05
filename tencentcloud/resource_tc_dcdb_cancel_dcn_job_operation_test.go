package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcdbCancelDcnJobOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbCancelDcnJobOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_cancel_dcn_job_operation.cancel_operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_cancel_dcn_job_operation.cancel_operation", "instance_id"),
				),
			},
		},
	})
}

const testAccDcdbHourInsOperation_dcn = defaultAzVariable + `
data "tencentcloud_security_groups" "internal" {
	name = "default"
}

data "tencentcloud_vpc_instances" "vpc" {
	name ="Default-VPC"
}
	
data "tencentcloud_vpc_subnets" "subnet" {
	vpc_id = data.tencentcloud_vpc_instances.vpc.instance_list.0.vpc_id
}

locals {
	vpc_id = data.tencentcloud_vpc_subnets.subnet.instance_list.0.vpc_id
	subnet_id = data.tencentcloud_vpc_subnets.subnet.instance_list.0.subnet_id
	sg_id = data.tencentcloud_security_groups.internal.security_groups.0.security_group_id
}

resource "tencentcloud_dcdb_hourdb_instance" "hourdb_instance_m" {
	instance_name = "test_dcdb_db_hourdb_instance_master"
	zones = [var.default_az]
	shard_memory = "2"
	shard_storage = "10"
	shard_node_count = "2"
	shard_count = "2"
	vpc_id = local.vpc_id
	subnet_id = local.subnet_id
	security_group_id = local.sg_id
	db_version_id = "8.0"
	dcn_region      = "ap-guangzhou"
	dcn_instance_id = local.dcn_dcdb_id
	resource_tags {
	  tag_key = "aaa"
	  tag_value = "bbb"
	}
}

resource "tencentcloud_dcdb_hourdb_instance" "hourdb_instance_d" {
	instance_name = "test_dcdb_db_hourdb_instance_dcn"
	zones = [var.default_az]
	shard_memory = "2"
	shard_storage = "10"
	shard_node_count = "2"
	shard_count = "2"
	vpc_id = local.vpc_id
	subnet_id = local.subnet_id
	security_group_id = local.sg_id
	db_version_id = "8.0"
	resource_tags {
	  tag_key = "aaa"
	  tag_value = "bbb"
	}
}

  locals {
	// master_dcdb_id = tencentcloud_dcdb_hourdb_instance.hourdb_instance_m.id
	dcn_dcdb_id = tencentcloud_dcdb_hourdb_instance.hourdb_instance_d.id
  }
`

const testAccDcdbCancelDcnJobOperation = testAccDcdbHourInsOperation_dcn + `

resource "tencentcloud_dcdb_cancel_dcn_job_operation" "cancel_operation" {
  instance_id = local.dcn_dcdb_id
}

`
