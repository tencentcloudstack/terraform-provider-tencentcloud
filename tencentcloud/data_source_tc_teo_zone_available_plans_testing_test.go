package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTestingTeoZoneAvailablePlansDataSource -v
func TestAccTencentCloudTestingTeoZoneAvailablePlansDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTestingTeoZoneAvailablePlans,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_teo_zone_available_plans.zone_available_plans"),
				),
			},
		},
	})
}

const testAccDataSourceTestingTeoZoneAvailablePlans = `

data "tencentcloud_teo_zone_available_plans" "zone_available_plans" {
}

`
