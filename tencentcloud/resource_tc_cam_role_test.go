package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCamRoleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCamRole,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cam_role.role", "id")),
			},
			{
				ResourceName:      "tencentcloud_cam_role.role",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCamRole = `

resource "tencentcloud_cam_role" "role" {
  role_name = ""
  policy_document = ""
  description = ""
  console_login = 
  session_duration = 
}

`
