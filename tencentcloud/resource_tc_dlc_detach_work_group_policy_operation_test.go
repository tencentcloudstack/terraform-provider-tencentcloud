package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDlcDetachWorkGroupPolicyOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcDetachWorkGroupPolicyOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_detach_work_group_policy_operation.detach_work_group_policy_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_dlc_detach_work_group_policy_operation.detach_work_group_policy_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDlcDetachWorkGroupPolicyOperation = `

resource "tencentcloud_dlc_detach_work_group_policy_operation" "detach_work_group_policy_operation" {
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
