package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRedisInstanceNodeInfoDataSource_basic -v
func TestAccTencentCloudRedisInstanceNodeInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstanceNodeInfoDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_redis_instance_node_info.instance_node_info"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_node_info.instance_node_info", "instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_node_info.instance_node_info", "proxy_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_node_info.instance_node_info", "proxy.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_node_info.instance_node_info", "proxy.0.node_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_node_info.instance_node_info", "proxy.0.zone_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_node_info.instance_node_info", "redis_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_node_info.instance_node_info", "redis.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_node_info.instance_node_info", "redis.0.node_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_node_info.instance_node_info", "redis.0.node_role"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_node_info.instance_node_info", "redis.0.zone_id"),
				),
			},
		},
	})
}

const testAccRedisInstanceNodeInfoDataSourceVar = `
variable "instance_id" {
	default = "` + defaultCrsInstanceId + `"
}
`
const testAccRedisInstanceNodeInfoDataSource = testAccRedisInstanceNodeInfoDataSourceVar + `

data "tencentcloud_redis_instance_node_info" "instance_node_info" {
	instance_id = var.instance_id
}

`
