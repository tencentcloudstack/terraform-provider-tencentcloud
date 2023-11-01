package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudCssStreamMonitorResource_basic -v
func TestAccTencentCloudCssStreamMonitorResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssStreamMonitor,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_stream_monitor.stream_monitor", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "ai_format_diagnose", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "ai_asr_input_index_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "ai_ocr_input_index_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "allow_monitor_report", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "asr_language", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "check_stream_broken", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "check_stream_low_frame_rate", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "monitor_name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "ocr_language", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "input_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "input_list.0.input_app", "live"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "input_list.0.input_domain", "177154.push.tlivecloud.com"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "input_list.0.input_stream_name", "ppp"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "notify_policy.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "notify_policy.0.callback_url", "http://example.com/test"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "notify_policy.0.notify_policy_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "output_info.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "output_info.0.output_domain", "test122.jingxhu.top"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "output_info.0.output_stream_height", "1080"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "output_info.0.output_stream_name", "afc7847d-1fe1-43bc-b1e4-20d86303c393"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "output_info.0.output_stream_width", "1920"),
				),
			},
			{
				ResourceName:      "tencentcloud_css_stream_monitor.stream_monitor",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCssStreamMonitorUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_stream_monitor.stream_monitor", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "ai_format_diagnose", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "ai_asr_input_index_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "ai_ocr_input_index_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "allow_monitor_report", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "asr_language", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "check_stream_broken", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "check_stream_low_frame_rate", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "monitor_name", "test1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "ocr_language", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "input_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "input_list.0.input_app", "live1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "input_list.0.input_domain", "177154.push.tlivecloud.com"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "input_list.0.input_stream_name", "pppq"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "notify_policy.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "notify_policy.0.callback_url", "http://example.com/test1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "notify_policy.0.notify_policy_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "output_info.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "output_info.0.output_domain", "test122.jingxhu.top"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "output_info.0.output_stream_height", "1080"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "output_info.0.output_stream_name", "afc7847d-1fe1-43bc-b1e4-20d86303c393"),
					resource.TestCheckResourceAttr("tencentcloud_css_stream_monitor.stream_monitor", "output_info.0.output_stream_width", "1920"),
				),
			},
		},
	})
}

const testAccCssStreamMonitor = `

resource "tencentcloud_css_stream_monitor" "stream_monitor" {
  ai_asr_input_index_list = [
    1,
  ]
  ai_format_diagnose = 1
  ai_ocr_input_index_list = [
    1,
  ]
  allow_monitor_report        = 1
  asr_language                = 1
  check_stream_broken         = 1
  check_stream_low_frame_rate = 1
  monitor_name                = "test"
  ocr_language                = 1

  input_list {
    input_app         = "live"
    input_domain      = "177154.push.tlivecloud.com"
    input_stream_name = "ppp"
  }

  notify_policy {
    callback_url       = "http://example.com/test"
    notify_policy_type = 1
  }

  output_info {
    output_domain        = "test122.jingxhu.top"
    output_stream_height = 1080
    output_stream_name   = "afc7847d-1fe1-43bc-b1e4-20d86303c393"
    output_stream_width  = 1920
  }
}

`

const testAccCssStreamMonitorUp = `
resource "tencentcloud_css_stream_monitor" "stream_monitor" {
  ai_asr_input_index_list = [
    1,
  ]
  ai_format_diagnose = 1
  ai_ocr_input_index_list = [
    1,
  ]
  allow_monitor_report        = 1
  asr_language                = 1
  check_stream_broken         = 1
  check_stream_low_frame_rate = 1
  monitor_name                = "test1"
  ocr_language                = 1

  input_list {
    input_app         = "live1"
    input_domain      = "177154.push.tlivecloud.com"
    input_stream_name = "pppq"
  }

  notify_policy {
    callback_url       = "http://example.com/test1"
    notify_policy_type = 1
  }

  output_info {
    output_domain        = "test122.jingxhu.top"
    output_stream_height = 1080
    output_stream_name   = "afc7847d-1fe1-43bc-b1e4-20d86303c393"
    output_stream_width  = 1920
  }
}
`
