package cynosdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCynosdbZoneConfig_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCynosdbZoneConfig(),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_zone_config.test"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone_config.test", "list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone_config.test", "list.0.cpu"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone_config.test", "list.0.memory"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone_config.test", "list.0.max_storage_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone_config.test", "list.0.min_storage_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone_config.test", "list.0.machine_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone_config.test", "list.0.max_io_bandwidth"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone_config.test", "list.0.zone_stock_infos.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone_config.test", "list.0.zone_stock_infos.0.zone"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_zone_config.test", "list.0.zone_stock_infos.0.has_stock"),
				),
			},
		},
	})
}

func testAccDataSourceCynosdbZoneConfig() string {
	return `data "tencentcloud_cynosdb_zone_config" "test" {
		
	}`
}
