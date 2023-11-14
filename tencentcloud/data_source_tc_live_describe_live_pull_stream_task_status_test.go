package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLiveDescribeLivePullStreamTaskStatusDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveDescribeLivePullStreamTaskStatusDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_live_describe_live_pull_stream_task_status.describe_live_pull_stream_task_status")),
			},
		},
	})
}

const testAccLiveDescribeLivePullStreamTaskStatusDataSource = `

data "tencentcloud_live_describe_live_pull_stream_task_status" "describe_live_pull_stream_task_status" {
  task_id = ""
  }

`
