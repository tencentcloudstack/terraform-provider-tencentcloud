package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapAccessRegionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapAccessRegionsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_access_regions.access_regions"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_access_regions.access_regions", "access_region_set.#"),
				),
			},
		},
	})
}

const testAccGaapAccessRegionsDataSource = `
data "tencentcloud_gaap_access_regions" "access_regions" {
}
`
