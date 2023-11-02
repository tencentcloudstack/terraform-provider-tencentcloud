package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixCssEnableOptimalSwitchingResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssEnableOptimalSwitching,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_enable_optimal_switching.enable_optimal_switching", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_enable_optimal_switching.enable_optimal_switching", "stream_name", "1308919341_test"),
					resource.TestCheckResourceAttr("tencentcloud_css_enable_optimal_switching.enable_optimal_switching", "enable_switch", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_css_enable_optimal_switching.enable_optimal_switching",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCssEnableOptimalSwitching = `

resource "tencentcloud_css_enable_optimal_switching" "enable_optimal_switching" {
  stream_name     = "1308919341_test"
  enable_switch   = 1
}

`
