package cdc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCdcSiteResource_basic -v
func TestAccTencentCloudNeedFixCdcSiteResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdcSite,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "country"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "province"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "city"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "address_line"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "optional_address_line"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "description"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "fiber_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "optical_standard"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "power_connectors"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "power_feed_drop"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "max_weight"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "power_draw_kva"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "uplink_speed_gbps"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "uplink_count"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "condition_requirement"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "dimension_requirement"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "redundant_networking"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "need_help"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "redundant_power"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "breaker_requirement"),
				),
			},
			{
				Config: testAccCdcSiteUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "country"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "province"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "city"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "address_line"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "optional_address_line"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "description"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "fiber_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "optical_standard"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "power_connectors"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "power_feed_drop"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "max_weight"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "power_draw_kva"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "uplink_speed_gbps"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "uplink_count"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "condition_requirement"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "dimension_requirement"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "redundant_networking"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "need_help"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "redundant_power"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_site.example", "breaker_requirement"),
				),
			},
			{
				ResourceName:      "tencentcloud_cdc_site.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdcSite = `
resource "tencentcloud_cdc_site" "example" {
  name                  = "tf-example"
  country               = "China"
  province              = "Guangdong Province"
  city                  = "Guangzhou"
  address_line          = "Shenzhen Tencent Building"
  optional_address_line = "Shenzhen Tencent Building of Binhai"
  description           = "desc."
  fiber_type            = "MM"
  optical_standard      = "MM"
  power_connectors      = "380VAC3P"
  power_feed_drop       = "DOWN"
  max_weight            = 100
  power_draw_kva        = 10
  uplink_speed_gbps     = 10
  uplink_count          = 2
  condition_requirement = true
  dimension_requirement = true
  redundant_networking  = true
  need_help             = true
  redundant_power       = true
  breaker_requirement   = true
}
`

const testAccCdcSiteUpdate = `
resource "tencentcloud_cdc_site" "example" {
  name                  = "tf-example-update"
  country               = "China"
  province              = "Guangdong Province"
  city                  = "Guangzhou"
  address_line          = "Shenzhen Tencent Building 001"
  optional_address_line = "Shenzhen Tencent Building of Binhai 002"
  description           = "desc update."
  fiber_type            = "MM"
  optical_standard      = "MM"
  power_connectors      = "380VAC3P"
  power_feed_drop       = "DOWN"
  max_weight            = 100
  power_draw_kva        = 10
  uplink_speed_gbps     = 10
  uplink_count          = 2
  condition_requirement = true
  dimension_requirement = true
  redundant_networking  = true
  need_help             = true
  redundant_power       = true
  breaker_requirement   = true
}
`
