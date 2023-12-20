package css_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudCssPullStreamTaskStatusDataSource_basic -v
func TestAccTencentCloudCssPullStreamTaskStatusDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssPullStreamTaskStatusDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_css_pull_stream_task_status.pull_stream_task_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_pull_stream_task_status.pull_stream_task_status", "task_status_info.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_pull_stream_task_status.pull_stream_task_status", "task_status_info.0.file_duration"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_pull_stream_task_status.pull_stream_task_status", "task_status_info.0.looped_times"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_pull_stream_task_status.pull_stream_task_status", "task_status_info.0.offset_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_pull_stream_task_status.pull_stream_task_status", "task_status_info.0.run_status"),
				),
			},
		},
	})
}

const testAccCssPullStreamTaskStatusDataSource = `

data "tencentcloud_css_pull_stream_task_status" "pull_stream_task_status" {
  task_id = "63229997"
}

`
