package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceRedisZoneConfig_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRedisZoneConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_redis_zone_config.test"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_zone_config.test", "list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_zone_config.test", "list.0.zone"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_zone_config.test", "list.0.type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_zone_config.test", "list.0.version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_zone_config.test", "list.0.mem_sizes.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_zone_config.test", "list.0.shard_memories.#"),
				),
			},
			{
				Config: testAccDataSourceRedisZoneConfigWithRegion(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_redis_zone_config.testWithRegion"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_zone_config.testWithRegion", "list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_zone_config.testWithRegion", "list.0.zone"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_zone_config.testWithRegion", "list.0.type_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_zone_config.testWithRegion", "list.0.version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_zone_config.testWithRegion", "list.0.mem_sizes.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_zone_config.testWithRegion", "list.0.shard_memories.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_zone_config.testWithRegion", "list.0.redis_shard_nums.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_zone_config.testWithRegion", "list.0.redis_replicas_nums.#"),
				),
			},
		},
	})
}

func testAccDataSourceRedisZoneConfig() string {
	return `data "tencentcloud_redis_zone_config" "test" {
		
	}`
}

func testAccDataSourceRedisZoneConfigWithRegion() string {
	return `data "tencentcloud_redis_zone_config" "testWithRegion" {
       region = "ap-guangzhou"
    }`
}
