package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudCssRestartPushTaskResource_basic -v
func TestAccTencentCloudCssRestartPushTaskResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssRestartPushTask,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_restart_push_task.restart_push_task", "id"),
				),
			},
		},
	})
}

const testAccCssRestartPushTask = testAccCssStreamMonitor + `

resource "tencentcloud_css_restart_push_task" "restart_push_task" {
  task_id = tencentcloud_css_stream_monitor.stream_monitor.id
  operator = "tf-test"
}

`
