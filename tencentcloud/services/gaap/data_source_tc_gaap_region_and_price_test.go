package gaap_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapRegionAndPriceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapRegionAndPriceDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_region_and_price.region_and_price"),
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
