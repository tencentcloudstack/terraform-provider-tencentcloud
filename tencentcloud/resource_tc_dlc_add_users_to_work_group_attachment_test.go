package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDlcAddUsersToWorkGroupAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcAddUsersToWorkGroupAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_add_users_to_work_group_attachment.add_users_to_work_group_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_dlc_add_users_to_work_group_attachment.add_users_to_work_group_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDlcAddUsersToWorkGroupAttachment = `

resource "tencentcloud_dlc_add_users_to_work_group_attachment" "add_users_to_work_group_attachment" {
  add_info {
		work_group_id = 23184
		user_ids = [100032676511]
  }
}

`
