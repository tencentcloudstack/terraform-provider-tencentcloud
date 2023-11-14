package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCamUserPermissionBoundaryResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCamUserPermissionBoundary,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cam_user_permission_boundary.user_permission_boundary", "id")),
			},
			{
				ResourceName:      "tencentcloud_cam_user_permission_boundary.user_permission_boundary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCamUserPermissionBoundary = `

resource "tencentcloud_cam_user_permission_boundary" "user_permission_boundary" {
  target_uin = 
  policy_id = 
}

`
