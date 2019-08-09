package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceRedisZoneConfig_basic(t *testing.T) {

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
				),
			},
			{
				Config: testAccDataSourceRedisZoneConfigWithRegion(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_redis_zone_config.testWithRegion"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_zone_config.testWithRegion", "list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_zone_config.testWithRegion", "list.0.zone"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_zone_config.testWithRegion", "list.0.type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_zone_config.testWithRegion", "list.0.version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_zone_config.testWithRegion", "list.0.mem_sizes.#"),
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
