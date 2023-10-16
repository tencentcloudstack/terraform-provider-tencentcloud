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
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cam_tag_role_attachment.tag_role", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_tag_role_attachment.tag_role", "tags.#"),
					resource.TestCheckResourceAttr("tencentcloud_cam_tag_role_attachment.tag_role", "role_id", "4611686018436805021")),
			},
			{
				ResourceName:      "tencentcloud_cam_tag_role_attachment.tag_role",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCamTagRole = `

resource "tencentcloud_cam_tag_role_attachment" "tag_role" {
  tags {
		key = "test1"
		value = "test1"

  }
  role_id = "4611686018436805021"
  }

`
