package dlc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcAddUsersToWorkGroupAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcAddUsersToWorkGroupAttachment,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_add_users_to_work_group_attachment.add_users_to_work_group_attachment", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_add_users_to_work_group_attachment.add_users_to_work_group_attachment", "add_info.#"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_add_users_to_work_group_attachment.add_users_to_work_group_attachment", "add_info.0.work_group_id", "23184"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_add_users_to_work_group_attachment.add_users_to_work_group_attachment", "add_info.0.user_ids.#"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_add_users_to_work_group_attachment.add_users_to_work_group_attachment", "add_info.0.user_ids.0", "100032676511"),
				),
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
