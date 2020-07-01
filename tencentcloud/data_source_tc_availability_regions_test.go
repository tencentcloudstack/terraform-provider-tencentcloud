package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudAvailabilityRegionsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudAvailabilityRegionsDataSourceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_availability_regions.all"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_availability_regions.all", "regions.#"),
				),
			},
			{
				Config: testAccTencentCloudAvailabilityRegionsDataSourceConfigFilterWithName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_availability_regions.filter"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_availability_regions.filter", "regions.#"),
				),
			},
			{
				Config: testAccTencentCloudAvailabilityRegionsDataSourceConfigIncludeUnavailable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_availability_regions.unavailable"),
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

const testAccTencentCloudAvailabilityRegionsDataSourceConfigFilterWithName = defaultVpcVariable + `
data "tencentcloud_availability_regions" "filter" {
  name = var.availability_region
}
`

const testAccTencentCloudAvailabilityRegionsDataSourceConfigIncludeUnavailable = `
data "tencentcloud_availability_regions" "unavailable" {
  include_unavailable = true
}
`
