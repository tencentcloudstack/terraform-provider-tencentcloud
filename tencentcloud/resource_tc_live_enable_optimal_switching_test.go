package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLiveEnableOptimalSwitchingResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveEnableOptimalSwitching,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_live_enable_optimal_switching.enable_optimal_switching", "id")),
			},
			{
				ResourceName:      "tencentcloud_live_enable_optimal_switching.enable_optimal_switching",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLiveEnableOptimalSwitching = `

resource "tencentcloud_live_enable_optimal_switching" "enable_optimal_switching" {
  stream_name = ""
  enable_switch = 
  host_group_name = ""
}

`
