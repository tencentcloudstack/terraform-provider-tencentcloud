package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceMysqlZoneConfig_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMysqlZoneConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_zone_config.test"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.test", "list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.test", "list.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.test", "list.0.is_default"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.test", "list.0.is_support_disaster_recovery"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.test", "list.0.is_support_vpc"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.test", "list.0.engine_versions.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.test", "list.0.pay_type.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.test", "list.0.support_slave_sync_modes.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.test", "list.0.disaster_recovery_zones.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.test", "list.0.slave_deploy_modes.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.test", "list.0.first_slave_zones.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.test", "list.0.second_slave_zones.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.test", "list.0.sells.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.test", "list.0.sells.0.mem_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.test", "list.0.sells.0.min_volume_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.test", "list.0.sells.0.max_volume_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.test", "list.0.sells.0.volume_step"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.test", "list.0.sells.0.qps"),
				),
			},
			{
				Config: testAccDataSourceMysqlZoneConfigWithRegion(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_zone_config.testWithRegion"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.testWithRegion", "list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.testWithRegion", "list.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.testWithRegion", "list.0.is_default"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.testWithRegion", "list.0.is_support_disaster_recovery"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.testWithRegion", "list.0.is_support_vpc"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.testWithRegion", "list.0.engine_versions.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.testWithRegion", "list.0.pay_type.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.testWithRegion", "list.0.support_slave_sync_modes.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.testWithRegion", "list.0.disaster_recovery_zones.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.testWithRegion", "list.0.slave_deploy_modes.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.testWithRegion", "list.0.first_slave_zones.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.testWithRegion", "list.0.second_slave_zones.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.testWithRegion", "list.0.sells.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.testWithRegion", "list.0.sells.0.mem_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.testWithRegion", "list.0.sells.0.min_volume_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.testWithRegion", "list.0.sells.0.max_volume_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.testWithRegion", "list.0.sells.0.volume_step"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_zone_config.testWithRegion", "list.0.sells.0.qps"),
				),
			},
		},
	})
}

func testAccDataSourceMysqlZoneConfig() string {
	return `data "tencentcloud_mysql_zone_config" "test" {
		
	}`
}

func testAccDataSourceMysqlZoneConfigWithRegion() string {
	return `data "tencentcloud_mysql_zone_config" "testWithRegion" {
       region = "ap-guangzhou"
    }`
}
