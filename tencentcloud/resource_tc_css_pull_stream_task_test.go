package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCssPullStreamTaskResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssPullStreamTask,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_css_pull_stream_task.pull_stream_task", "id")),
			},
			{
				ResourceName:      "tencentcloud_css_pull_stream_task.pull_stream_task",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCssPullStreamTask = `

resource "tencentcloud_css_pull_stream_task" "pull_stream_task" {
  source_type = "PullLivePushLive"
  source_urls = &lt;nil&gt;
  domain_name = &lt;nil&gt;
  app_name = &lt;nil&gt;
  stream_name = &lt;nil&gt;
  start_time = &lt;nil&gt;
  end_time = &lt;nil&gt;
  operator = &lt;nil&gt;
  push_args = &lt;nil&gt;
  callback_events = &lt;nil&gt;
  vod_loop_times = &lt;nil&gt;
  vod_refresh_type = &lt;nil&gt;
  callback_url = &lt;nil&gt;
  extra_cmd = &lt;nil&gt;
  comment = &lt;nil&gt;
  to_url = &lt;nil&gt;
  backup_source_type = &lt;nil&gt;
  backup_source_url = &lt;nil&gt;
  watermark_list {
		picture_url = &lt;nil&gt;
		x_position = &lt;nil&gt;
		y_position = &lt;nil&gt;
		width = &lt;nil&gt;
		height = &lt;nil&gt;
		location = &lt;nil&gt;

  }
  status = &lt;nil&gt;
          file_index = &lt;nil&gt;
  offset_time = &lt;nil&gt;
  }

`
