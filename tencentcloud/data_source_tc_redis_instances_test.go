package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudRedisInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudRedisInstancesDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instances.redis", "instance_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instances.redis", "instance_list.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instances.redis", "instance_list.0.zone"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instances.redis", "instance_list.0.project_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instances.redis", "instance_list.0.type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instances.redis", "instance_list.0.mem_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instances.redis", "instance_list.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instances.redis", "instance_list.0.ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instances.redis", "instance_list.0.port"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instances.redis", "instance_list.0.create_time"),
				),
			},
		},
	})
}

func testAccTencentCloudRedisInstancesDataSourceConfig() string {
	return `

resource "tencentcloud_redis_instance" "redis_instance_test"{
	availability_zone="ap-guangzhou-3"
	type="master_slave_redis"
	password="test12345789"
	mem_size=8192
	name="terrform_test"
	port=6379
}

data "tencentcloud_redis_instances" "redis" {
	zone = "ap-guangzhou-3"
	search_key = "${tencentcloud_redis_instance.redis_instance_test.id}"
}
	`
}
