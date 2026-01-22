package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataWorkflowPermissionsResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataWorkflowPermissions,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_workflow_permissions.example", "id"),
				),
			},
			{
				Config: testAccWedataWorkflowPermissionsUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_workflow_permissions.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_workflow_permissions.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWedataWorkflowPermissions = `
resource "tencentcloud_wedata_workflow_permissions" "example" {
  project_id  = "3108707295180644352"
  entity_id   = "53e78f97-f145-11f0-ba36-b8cef6a5af5c"
  entity_type = "folder"
  permission_list {
    permission_target_type = "user"
    permission_target_id   = "100028448903"
    permission_type_list   = ["CAN_MANAGE"]
  }

  permission_list {
    permission_target_type = "role"
    permission_target_id   = "308335260676890624"
    permission_type_list   = ["CAN_MANAGE"]
  }
}
`

const testAccWedataWorkflowPermissionsUpdate = `
resource "tencentcloud_wedata_workflow_permissions" "example" {
  project_id  = "3108707295180644352"
  entity_id   = "53e78f97-f145-11f0-ba36-b8cef6a5af5c"
  entity_type = "folder"

  permission_list {
    permission_target_type = "role"
    permission_target_id   = "308335260676890624"
    permission_type_list   = ["CAN_MANAGE"]
  }
}
`
