package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcAttachWorkGroupPolicyOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcAttachWorkGroupPolicyOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_attach_work_group_policy_operation.attach_work_group_policy_operation", "id")),
			},
		},
	})
}

const testAccDlcAttachWorkGroupPolicyOperation = `

resource "tencentcloud_dlc_attach_work_group_policy_operation" "attach_work_group_policy_operation" {
  work_group_id = 23184
  policy_set {
		database = "test_iac_keep"
		catalog = "DataLakeCatalog"
		table = ""
		operation = "ASSAYER"
		policy_type = "DATABASE"
		function = ""
		view = ""
		column = ""
		source = "USER"
		mode = "COMMON"
  }
}

`
