package tencentcloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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

					resource.TestMatchResourceAttr("data.tencentcloud_redis_instances.redis-tags", "instance_list.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instances.redis-tags", "instance_list.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instances.redis-tags", "instance_list.0.zone"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instances.redis-tags", "instance_list.0.project_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instances.redis-tags", "instance_list.0.type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instances.redis-tags", "instance_list.0.mem_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instances.redis-tags", "instance_list.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instances.redis-tags", "instance_list.0.ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instances.redis-tags", "instance_list.0.port"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instances.redis-tags", "instance_list.0.create_time"),
					resource.TestCheckResourceAttr("data.tencentcloud_redis_instances.redis-tags", "instance_list.0.tags.test", "test"),
				),
			},
		},
	})
}

func testAccTencentCloudRedisInstancesDataSourceConfig() string {
	return `
resource "tencentcloud_redis_instance" "redis_instance_test" {
  availability_zone = "ap-guangzhou-3"
  type_id           = 2
  password          = "test12345789"
  mem_size          = 8192
  name              = "terraform_test"
  port              = 6379

  tags = {
    "test" = "test"
  }
}

data "tencentcloud_redis_instances" "redis" {
  zone       = "ap-guangzhou-3"
  search_key = tencentcloud_redis_instance.redis_instance_test.id
}

data "tencentcloud_redis_instances" "redis-tags" {
  zone = "ap-guangzhou-3"
  tags = tencentcloud_redis_instance.redis_instance_test.tags
}
`
}
