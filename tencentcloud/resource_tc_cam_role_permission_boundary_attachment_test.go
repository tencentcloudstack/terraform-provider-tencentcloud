package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCamRolePermissionBoundaryAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCamRolePermissionBoundaryAttachment,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cam_role_permission_boundary_attachment.role_permission_boundary_attachment", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cam_role_permission_boundary_attachment.role_permission_boundary_attachment", "policy_id", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cam_role_permission_boundary_attachment.role_permission_boundary_attachment", "role_name", "test-cam-tag")),
			},
			{
				ResourceName:      "tencentcloud_cam_role_permission_boundary_attachment.role_permission_boundary_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCamRolePermissionBoundaryAttachment = `

resource "tencentcloud_cam_role_permission_boundary_attachment" "role_permission_boundary_attachment" {
  policy_id = 1
  role_name = "test-cam-tag"
}

`
