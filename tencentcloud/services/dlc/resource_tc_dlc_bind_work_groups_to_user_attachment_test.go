package dlc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcBindWorkGroupsToUserAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcBindWorkGroupsToUser,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_bind_work_groups_to_user_attachment.bind_work_groups_to_user", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_bind_work_groups_to_user_attachment.bind_work_groups_to_user", "add_info.#"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_bind_work_groups_to_user_attachment.bind_work_groups_to_user", "add_info.0.user_id", "100032772113"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_bind_work_groups_to_user_attachment.bind_work_groups_to_user", "add_info.0.work_group_ids.#"),
				),
			},
			{
				ResourceName:      "tencentcloud_dlc_bind_work_groups_to_user_attachment.bind_work_groups_to_user",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDlcBindWorkGroupsToUser = `

resource "tencentcloud_dlc_bind_work_groups_to_user_attachment" "bind_work_groups_to_user" {
  add_info {
    user_id = "100032772113"
    work_group_ids = [23184,23181]
  }
}

`
