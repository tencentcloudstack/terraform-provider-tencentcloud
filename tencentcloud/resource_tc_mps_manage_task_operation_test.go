package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsManageTaskOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsManageTaskOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_manage_task_operation.manage_task_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_manage_task_operation.manage_task_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsManageTaskOperation = `

resource "tencentcloud_mps_manage_task_operation" "manage_task_operation" {
  operation_type = ""
  task_id = ""
}

`
