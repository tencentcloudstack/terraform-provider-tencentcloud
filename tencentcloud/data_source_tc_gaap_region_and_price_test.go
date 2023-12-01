package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapRegionAndPriceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapRegionAndPriceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_region_and_price.region_and_price"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_region_and_price.region_and_price", "dest_region_set.#"),
				),
			},
		},
	})
}

const testAccGaapRegionAndPriceDataSource = `

data "tencentcloud_gaap_region_and_price" "region_and_price" {
}
`
