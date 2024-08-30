package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudIdentityCenterUserGroupAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityCenterUserGroupAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_user_group_attachment.identity_center_user_group_attachment", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_user_group_attachment.identity_center_user_group_attachment", "user_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_user_group_attachment.identity_center_user_group_attachment", "group_id"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_user_group_attachment.identity_center_user_group_attachment", "zone_id", "z-s64jh54hbcra"),
				),
			},
			{
				ResourceName:      "tencentcloud_identity_center_user_group_attachment.identity_center_user_group_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccIdentityCenterUserGroupAttachment = `
resource "tencentcloud_identity_center_group" "identity_center_group" {
	zone_id = "z-s64jh54hbcra"
    group_name = "attachment-group"
    description = "test"
}

resource "tencentcloud_identity_center_user" "identity_center_user" {
    zone_id = "z-s64jh54hbcra"
    user_name = "attachment-user"
    description = "test"
}

resource "tencentcloud_identity_center_user_group_attachment" "identity_center_user_group_attachment" {
	zone_id = "z-s64jh54hbcra"
	user_id = tencentcloud_identity_center_user.identity_center_user.user_id
	group_id = tencentcloud_identity_center_group.identity_center_group.group_id
}
`
