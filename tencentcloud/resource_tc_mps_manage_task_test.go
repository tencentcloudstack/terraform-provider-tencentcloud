package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsManageTaskResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsManageTask,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_manage_task.manage_task", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_manage_task.manage_task",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsManageTask = `

resource "tencentcloud_mps_manage_task" "manage_task" {
  operation_type = ""
  task_id = ""
}

`
