package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudPostgresqlRegionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlRegionsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_postgresql_regions.regions"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_regions.regions", "region_set.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_regions.regions", "region_set.0.region"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_regions.regions", "region_set.0.region_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_regions.regions", "region_set.0.region_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_regions.regions", "region_set.0.region_state"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_regions.regions", "region_set.0.support_international"),
				),
			},
		},
	})
}

const testAccPostgresqlRegionsDataSource = `

data "tencentcloud_postgresql_regions" "regions" {}

`
