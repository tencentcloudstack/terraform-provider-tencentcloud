package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTeoZoneAvailablePlansDataSource -v
func TestAccTencentCloudTeoZoneAvailablePlansDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTeoZoneAvailablePlans,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_teo_zone_available_plans.zone_available_plans"),
				),
			},
		},
	})
}

const testAccDataSourceTeoZoneAvailablePlans = `

data "tencentcloud_teo_zone_available_plans" "zone_available_plans" {
}

`
