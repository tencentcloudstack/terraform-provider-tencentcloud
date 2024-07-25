package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdcSiteResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdcSite,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.site", "id")),
			},
			{
				ResourceName:      "tencentcloud_cdc_site.site",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdcSite = `

resource "tencentcloud_cdc_site" "site" {
  name = "site-1"
  country = "China"
  province = "Guangdong Province"
  city = "Guangzhou"
  address_line = ""
  description = ""
  note = ""
  fiber_type = ""
  optical_standard = ""
  power_connectors = ""
  power_feed_drop = ""
  max_weight = 100
  power_draw_kva = 10
  uplink_speed_gbps = 10
  uplink_count = 2
  condition_requirement = true
  dimension_requirement = true
  redundant_networking = true
  postal_code = 10000
  optional_address_line = ""
  need_help = true
  redundant_power = true
  breaker_requirement = true
}

`
