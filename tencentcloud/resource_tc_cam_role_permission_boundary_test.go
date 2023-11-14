package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCamRolePermissionBoundaryResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCamRolePermissionBoundary,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cam_role_permission_boundary.role_permission_boundary", "id")),
			},
			{
				ResourceName:      "tencentcloud_cam_role_permission_boundary.role_permission_boundary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCamRolePermissionBoundary = `

resource "tencentcloud_cam_role_permission_boundary" "role_permission_boundary" {
  policy_id = 
  role_id = ""
  role_name = ""
}

`
