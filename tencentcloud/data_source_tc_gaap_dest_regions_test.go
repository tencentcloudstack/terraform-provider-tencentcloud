package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapDestRegionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapDestRegionsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_dest_regions.dest_regions"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_dest_regions.dest_regions", "dest_region_set.#"),
				),
			},
		},
	})
}

const testAccGaapDestRegionsDataSource = `
data "tencentcloud_gaap_dest_regions" "dest_regions" {
}
`
