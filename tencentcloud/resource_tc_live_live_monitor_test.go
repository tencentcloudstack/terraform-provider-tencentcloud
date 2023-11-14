package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLiveLive_monitorResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveLive_monitor,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_live_live_monitor.live_monitor", "id")),
			},
			{
				ResourceName:      "tencentcloud_live_live_monitor.live_monitor",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLiveLive_monitor = `

resource "tencentcloud_live_live_monitor" "live_monitor" {
  monitor_id = ""
  audible_input_index_list = 
}

`
