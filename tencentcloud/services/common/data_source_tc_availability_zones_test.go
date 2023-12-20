package common_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAvailabilityZonesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudAvailabilityZonesDataSourceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_availability_zones.all"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_availability_zones.all", "zones.#"),
				),
			},
			{
				Config: testAccTencentCloudAvailabilityZonesDataSourceConfigFilterWithName,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_availability_zones.filter"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_availability_zones.filter", "zones.#"),
				),
			},
			{
				Config: testAccTencentCloudAvailabilityZonesDataSourceConfigIncludeUnavailable,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_availability_zones.unavailable"),
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

const testAccTencentCloudAvailabilityZonesDataSourceConfigFilterWithName = tcacctest.DefaultVpcVariable + `
data "tencentcloud_availability_zones" "filter" {
  name = var.availability_zone
}
`

const testAccTencentCloudAvailabilityZonesDataSourceConfigIncludeUnavailable = `
data "tencentcloud_availability_zones" "unavailable" {
  include_unavailable = true
}
`
