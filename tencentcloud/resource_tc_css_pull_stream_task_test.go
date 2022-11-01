package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCssPullStreamTask_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssPullStreamTask,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_pull_stream_task.pull_stream_task", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_css_pull_stream_task.pullStreamTask",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCssPullStreamTask = `

resource "tencentcloud_css_pull_stream_task" "pull_stream_task" {
  source_type = "PullLivePushLive"
  source_urls = ""
  domain_name = ""
  app_name = ""
  stream_name = ""
  start_time = ""
  end_time = ""
  operator = ""
  push_args = ""
  callback_events = ""
  vod_loop_times = ""
  vod_refresh_type = ""
  callback_url = ""
  extra_cmd = ""
  comment = ""
  to_url = ""
  backup_source_type = ""
  backup_source_url = ""
  watermark_list {
			picture_url = ""
			x_position = ""
			y_position = ""
			width = ""
			height = ""
			location = ""

  }
  status = ""
          file_index = ""
  offset_time = ""
  }

`
