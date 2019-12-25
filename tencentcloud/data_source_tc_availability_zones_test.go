package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudAvailabilityZonesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudAvailabilityZonesDataSourceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_availability_zones.all"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_availability_zones.all", "zones.#"),
				),
			},
			{
				Config: testAccTencentCloudAvailabilityZonesDataSourceConfigFilterWithName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_availability_zones.filter"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_availability_zones.filter", "zones.#"),
				),
			},
			{
				Config: testAccTencentCloudAvailabilityZonesDataSourceConfigIncludeUnavailable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_availability_zones.unavailable"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_availability_zones.unavailable", "zones.#"),
				),
			},
		},
	})
}

const testAccTencentCloudAvailabilityZonesDataSourceConfigBasic = `
data "tencentcloud_availability_zones" "all" {
}
`

const testAccTencentCloudAvailabilityZonesDataSourceConfigFilterWithName = defaultVpcVariable + `
data "tencentcloud_availability_zones" "filter" {
  name = var.availability_zone
}
`

const testAccTencentCloudAvailabilityZonesDataSourceConfigIncludeUnavailable = `
data "tencentcloud_availability_zones" "unavailable" {
  include_unavailable = true
}
`
