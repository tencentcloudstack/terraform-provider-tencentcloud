package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLiveStreamMonitorResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveStreamMonitor,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_live_stream_monitor.stream_monitor", "id")),
			},
			{
				ResourceName:      "tencentcloud_live_stream_monitor.stream_monitor",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLiveStreamMonitor = `

resource "tencentcloud_live_stream_monitor" "stream_monitor" {
  output_info {
		output_stream_width = 
		output_stream_height = 
		output_stream_name = ""
		output_domain = ""
		output_app = ""

  }
  input_list {
		input_stream_name = ""
		input_domain = ""
		input_app = ""
		input_url = ""
		description = ""

  }
  monitor_name = ""
  notify_policy {
		notify_policy_type = 
		callback_url = ""

  }
  asr_language = 
  ocr_language = 
  ai_asr_input_index_list = 
  ai_ocr_input_index_list = 
  check_stream_broken = 
  check_stream_low_frame_rate = 
  allow_monitor_report = 
  ai_format_diagnose = 
}

`
