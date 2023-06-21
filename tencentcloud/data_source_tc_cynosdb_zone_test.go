package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbZoneDataSource_basic -v
func TestAccTencentCloudCynosdbZoneDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbZoneDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_zone.zone"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone.zone", "region_set.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone.zone", "region_set.0.db_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone.zone", "region_set.0.modules.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone.zone", "region_set.0.modules.0.is_disable"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone.zone", "region_set.0.modules.0.module_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone.zone", "region_set.0.region"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone.zone", "region_set.0.region_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone.zone", "region_set.0.region_zh"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone.zone", "region_set.0.zone_set.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone.zone", "region_set.0.zone_set.0.has_permission"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone.zone", "region_set.0.zone_set.0.is_support_normal"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone.zone", "region_set.0.zone_set.0.is_support_serverless"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone.zone", "region_set.0.zone_set.0.is_whole_rdma_zone"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone.zone", "region_set.0.zone_set.0.physical_zone"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone.zone", "region_set.0.zone_set.0.zone"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone.zone", "region_set.0.zone_set.0.zone_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone.zone", "region_set.0.zone_set.0.zone_zh"),
				),
			},
		},
	})
}

const testAccCynosdbZoneDataSource = `

data "tencentcloud_cynosdb_zone" "zone" {
  include_virtual_zones = true
  show_permission = true
}

`
