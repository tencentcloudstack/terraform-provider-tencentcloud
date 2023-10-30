package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDlcAttachWorkGroupPolicyAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcAttachWorkGroupPolicyAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_attach_work_group_policy_attachment.attach_work_group_policy_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_dlc_attach_work_group_policy_attachment.attach_work_group_policy_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDlcAttachWorkGroupPolicyAttachment = `

resource "tencentcloud_dlc_attach_work_group_policy_attachment" "attach_work_group_policy_attachment" {
  work_group_id = 122
  policy_set {
		database = "*"
		catalog = "*"
		table = "*"
		operation = "ALL"
		policy_type = "ADMIN"
		function = "*"
		view = "*"
		column = "*"
		data_engine = "*"
		re_auth = false
		source = "USER"
		mode = "COMMON"
		operator = "admin"
		create_time = ""
		source_id = 
		source_name = ""
		id = 1

  }
}

`
