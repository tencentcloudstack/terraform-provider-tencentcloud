package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLivePullpush_taskResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLivePullpush_task,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_live_pullpush_task.pullpush_task", "id")),
			},
			{
				ResourceName:      "tencentcloud_live_pullpush_task.pullpush_task",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLivePullpush_task = `

resource "tencentcloud_live_pullpush_task" "pullpush_task" {
  task_id = ""
  operator = ""
}

`
