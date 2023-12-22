package postgresql_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudPostgresqlZonesDataSource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccStepSetRegion(t, "ap-guangzhou")
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlZonesDataSource,
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-guangzhou")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_postgresql_zones.zones"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_zones.zones", "zone_set.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_zones.zones", "zone_set.0.zone"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_zones.zones", "zone_set.0.zone_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_zones.zones", "zone_set.0.zone_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_zones.zones", "zone_set.0.zone_state"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_zones.zones", "zone_set.0.zone_support_ipv6"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_zones.zones", "zone_set.0.standby_zone_set.#"),
				),
			},
		},
	})
}

const testAccPostgresqlZonesDataSource = `

data "tencentcloud_postgresql_zones" "zones" {}

`
