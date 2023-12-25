package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudTestingTeoZoneAvailablePlansDataSource -v
func TestAccTencentCloudTestingTeoZoneAvailablePlansDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTestingTeoZoneAvailablePlans,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_zone_available_plans.zone_available_plans"),
				),
			},
		},
	})
}

const testAccDataSourceTestingTeoZoneAvailablePlans = `

data "tencentcloud_teo_zone_available_plans" "zone_available_plans" {
}

`
