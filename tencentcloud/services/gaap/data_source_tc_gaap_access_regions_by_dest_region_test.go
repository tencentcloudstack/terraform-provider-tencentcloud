package gaap_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapAccessRegionsByDestRegionDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapAccessRegionsByDestRegionDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_access_regions_by_dest_region.access_regions_by_dest_region"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_access_regions_by_dest_region.access_regions_by_dest_region", "access_region_set.#"),
				),
			},
		},
	})
}

const testAccGaapAccessRegionsByDestRegionDataSource = `
data "tencentcloud_gaap_access_regions_by_dest_region" "access_regions_by_dest_region" {
	dest_region = "SouthChina"
}
`
