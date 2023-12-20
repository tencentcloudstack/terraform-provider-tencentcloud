package common_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAvailabilityRegionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudAvailabilityRegionsDataSourceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_availability_regions.all"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_availability_regions.all", "regions.#"),
				),
			},
			{
				Config: testAccTencentCloudAvailabilityRegionsDataSourceConfigFilterWithName,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_availability_regions.filter"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_availability_regions.filter", "regions.#"),
				),
			},
			{
				Config: testAccTencentCloudAvailabilityRegionsDataSourceConfigIncludeUnavailable,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_availability_regions.unavailable"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_availability_regions.unavailable", "regions.#"),
				),
			},
		},
	})
}

const testAccTencentCloudAvailabilityRegionsDataSourceConfigBasic = `
data "tencentcloud_availability_regions" "all" {
}
`

const testAccTencentCloudAvailabilityRegionsDataSourceConfigFilterWithName = tcacctest.DefaultVpcVariable + `
data "tencentcloud_availability_regions" "filter" {
  name = var.availability_region
}
`

const testAccTencentCloudAvailabilityRegionsDataSourceConfigIncludeUnavailable = `
data "tencentcloud_availability_regions" "unavailable" {
  include_unavailable = true
}
`
