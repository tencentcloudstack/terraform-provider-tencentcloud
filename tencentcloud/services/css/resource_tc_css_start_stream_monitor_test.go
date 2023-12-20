package css_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCssStartStreamMonitorResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssStartStreamMonitor,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_start_stream_monitor.start_stream_monitor", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_css_start_stream_monitor.start_stream_monitor",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCssStartStreamMonitor = testAccCssStreamMonitor + `

resource "tencentcloud_css_start_stream_monitor" "start_stream_monitor" {
  monitor_id = tencentcloud_css_stream_monitor.stream_monitor.id
  audible_input_index_list = [1]
}

`
