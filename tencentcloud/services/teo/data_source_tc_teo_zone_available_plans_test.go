package teo_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTeoZoneAvailablePlansDataSource -v
func TestAccTencentCloudTeoZoneAvailablePlansDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTeoZoneAvailablePlans,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_zone_available_plans.zone_available_plans"),
				),
			},
		},
	})
}

const testAccDataSourceTeoZoneAvailablePlans = `

data "tencentcloud_teo_zone_available_plans" "zone_available_plans" {
}

`
