package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCamTagRoleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCamTagRole,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cam_tag_role.tag_role", "id")),
			},
			{
				ResourceName:      "tencentcloud_cam_tag_role.tag_role",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCamTagRole = `

resource "tencentcloud_cam_tag_role" "tag_role" {
  role_name = ""
  role_id = ""
}

`
