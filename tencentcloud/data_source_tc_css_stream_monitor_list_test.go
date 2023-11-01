package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudCssStreamMonitorListDataSource_basic -v
func TestAccTencentCloudCssStreamMonitorListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssStreamMonitorListDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_css_stream_monitor_list.stream_monitor_list"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.ai_asr_input_index_list.#"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.ai_format_diagnose.#"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.ai_ocr_input_index_list.#"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.allow_monitor_report.#"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.asr_language.#"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.audible_input_index_list.#"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.check_stream_broken.#"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.check_stream_low_frame_rate.#"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.create_time.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.input_list.#"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.input_list.0.input_app"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.input_list.0.input_domain"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.input_list.0.input_stream_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.monitor_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.monitor_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.notify_policy.#"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.notify_policy.0.notify_policy_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.ocr_language"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.output_info.#"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.output_info.0.output_domain"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.output_info.0.output_stream_height"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.output_info.0.output_stream_name"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.output_info.0.output_stream_width"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.start_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.stop_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.update_time"),
				),
			},
		},
	})
}

const testAccCssStreamMonitorListDataSource = `

data "tencentcloud_css_stream_monitor_list" "stream_monitor_list" {
}

`
