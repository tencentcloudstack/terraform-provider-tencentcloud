package gaap_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapAccessRegionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapAccessRegionsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_access_regions.access_regions"),
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
