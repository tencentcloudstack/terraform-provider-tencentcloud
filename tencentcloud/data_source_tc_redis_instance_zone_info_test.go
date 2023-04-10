package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRedisInstanceZoneInfoDataSource_basic -v
func TestAccTencentCloudRedisInstanceZoneInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstanceZoneInfoDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_redis_instance_zone_info.instance_zone_info"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_zone_info.instance_zone_info", "replica_groups.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_zone_info.instance_zone_info", "replica_groups.0.group_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_zone_info.instance_zone_info", "replica_groups.0.group_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_zone_info.instance_zone_info", "replica_groups.0.redis_nodes.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_zone_info.instance_zone_info", "replica_groups.0.redis_nodes.0.keys"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_zone_info.instance_zone_info", "replica_groups.0.redis_nodes.0.node_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_zone_info.instance_zone_info", "replica_groups.0.redis_nodes.0.role"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_zone_info.instance_zone_info", "replica_groups.0.redis_nodes.0.slot"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_zone_info.instance_zone_info", "replica_groups.0.redis_nodes.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_zone_info.instance_zone_info", "replica_groups.0.role"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_zone_info.instance_zone_info", "replica_groups.0.zone_id"),
				),
			},
		},
	})
}

const testAccRedisInstanceZoneInfoDataSourceVar = `
variable "vpc_id" {
	default = "` + defaultCrsVpcId + `"
}
variable "subnet_id" {
	default = "` + defaultCrsSubnetId + `"
}
`

const testAccRedisInstanceZoneInfoDataSource = testAccRedisInstanceZoneInfoDataSourceVar + `

resource "tencentcloud_redis_instance" "redis_instance_test" {
    auto_renew_flag    = 0
    availability_zone  = "ap-guangzhou-6"
    charge_type        = "POSTPAID"
    mem_size           = 4096
    name               = "terraform-test"
    no_auth            = true
    port               = 6379
    project_id         = 0
    redis_replicas_num = 1
    redis_shard_num    = 1
    replica_zone_ids   = [
        100007,
    ]
    replicas_read_only = false
    security_groups    = [
        "sg-ijato2x1",
    ]
    tags               = {}
    type_id            = 6
	vpc_id 			   = var.vpc_id
  	subnet_id		   = var.subnet_id
}

data "tencentcloud_redis_instance_zone_info" "instance_zone_info" {
  instance_id = tencentcloud_redis_instance.redis_instance_test.id
}

`
