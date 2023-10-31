package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDlcDetachUserPolicyOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcDetachUserPolicyOperation,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_detach_user_policy_operation.detach_user_policy_operation", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_detach_user_policy_operation.detach_user_policy_operation", "user_id", "100032676511"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_detach_user_policy_operation.detach_user_policy_operation", "policy_set.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_detach_user_policy_operation.detach_user_policy_operation", "policy_set.0.database"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_detach_user_policy_operation.detach_user_policy_operation", "policy_set.0.catalog"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_detach_user_policy_operation.detach_user_policy_operation", "policy_set.0.table"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_detach_user_policy_operation.detach_user_policy_operation", "policy_set.0.operation"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_detach_user_policy_operation.detach_user_policy_operation", "policy_set.0.policy_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_detach_user_policy_operation.detach_user_policy_operation", "policy_set.0.re_auth"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_detach_user_policy_operation.detach_user_policy_operation", "policy_set.0.source"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_detach_user_policy_operation.detach_user_policy_operation", "policy_set.0.mode"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_detach_user_policy_operation.detach_user_policy_operation", "policy_set.0.operator"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_detach_user_policy_operation.detach_user_policy_operation", "policy_set.0.id"),
				),
			},
		},
	})
}

const testAccDlcDetachUserPolicyOperation = `

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
data "tencentcloud_dlc_describe_user_info" "describe_user_info" {
  user_id = "100032676511"
  type = "DataAuth"
  sort_by = "create-time"
  sorting = "desc"
  depends_on = [
    tencentcloud_dlc_attach_user_policy_operation.attach_user_policy_operation,
  ]
}
resource "tencentcloud_dlc_detach_user_policy_operation" "detach_user_policy_operation" {
  user_id = "100032676511"
  policy_set {
    database = data.tencentcloud_dlc_describe_user_info.describe_user_info.user_info.0.data_policy_info.0.policy_set.0.database
    catalog = data.tencentcloud_dlc_describe_user_info.describe_user_info.user_info.0.data_policy_info.0.policy_set.0.catalog
    table = data.tencentcloud_dlc_describe_user_info.describe_user_info.user_info.0.data_policy_info.0.policy_set.0.table
    operation = data.tencentcloud_dlc_describe_user_info.describe_user_info.user_info.0.data_policy_info.0.policy_set.0.operation
    policy_type = data.tencentcloud_dlc_describe_user_info.describe_user_info.user_info.0.data_policy_info.0.policy_set.0.policy_type
    re_auth = data.tencentcloud_dlc_describe_user_info.describe_user_info.user_info.0.data_policy_info.0.policy_set.0.re_auth
    source = data.tencentcloud_dlc_describe_user_info.describe_user_info.user_info.0.data_policy_info.0.policy_set.0.source
    mode = data.tencentcloud_dlc_describe_user_info.describe_user_info.user_info.0.data_policy_info.0.policy_set.0.mode
    operator = data.tencentcloud_dlc_describe_user_info.describe_user_info.user_info.0.data_policy_info.0.policy_set.0.operator
    id = data.tencentcloud_dlc_describe_user_info.describe_user_info.user_info.0.data_policy_info.0.policy_set.0.id
  }
  depends_on = [
    tencentcloud_dlc_attach_user_policy_operation.attach_user_policy_operation,
  ]
}

`
