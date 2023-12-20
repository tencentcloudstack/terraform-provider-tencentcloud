package css_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudCssStreamMonitorListDataSource_basic -v
func TestAccTencentCloudCssStreamMonitorListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssStreamMonitorListDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_css_stream_monitor_list.stream_monitor_list"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.ai_asr_input_index_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.input_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.monitor_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.monitor_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.notify_policy.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.ocr_language"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_stream_monitor_list.stream_monitor_list", "live_stream_monitors.0.output_info.#"),
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
