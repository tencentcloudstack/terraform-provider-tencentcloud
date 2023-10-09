package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCamUserPermissionBoundaryAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCamUserPermissionBoundary,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cam_user_permission_boundary_attachment.user_permission_boundary", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user_permission_boundary_attachment.user_permission_boundary", "target_uin", "100032767426"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user_permission_boundary_attachment.user_permission_boundary", "policy_id", "151113272"),
				),
			},
			{
				ResourceName:            "tencentcloud_cam_user_permission_boundary_attachment.user_permission_boundary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"target_uin"},
			},
		},
	})
}

const testAccCamUserPermissionBoundary = `

resource "tencentcloud_cam_user_permission_boundary_attachment" "user_permission_boundary" {
  target_uin = 100032767426
  policy_id = 151113272
}

`
