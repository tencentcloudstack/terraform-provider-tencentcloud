package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTencentCloudTeoZoneShieldArea(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoZoneShieldArea,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tencentcloud_teo_ddos_policy.foo", "app_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_teo_ddos_policy.foo", "shield_areas"),
				),
			},
		},
	})
}

const testAccTeoZoneShieldArea = `
data "tencentcloud_teo_ddos_policy" "foo" {
}
`
