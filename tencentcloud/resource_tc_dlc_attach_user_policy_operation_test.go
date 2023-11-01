package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcAttachUserPolicyOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcAttachUserPolicyOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_attach_user_policy_operation.attach_user_policy_operation", "id")),
			},
		},
	})
}

const testAccDlcAttachUserPolicyOperation = `

resource "tencentcloud_dlc_attach_user_policy_operation" "attach_user_policy_operation" {
  user_id = "100032676511"
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
