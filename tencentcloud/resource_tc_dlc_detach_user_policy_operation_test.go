package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
					resource.TestCheckResourceAttr("tencentcloud_dlc_detach_user_policy_operation.detach_user_policy_operation", "policy_set.0.table", ""),
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

resource "tencentcloud_dlc_detach_user_policy_operation" "detach_user_policy_operation" {
 user_id = 100032676511
 policy_set {
   database = "test_iac_keep"
   catalog = "DataLakeCatalog"
   table = ""
   operation = "ASSAYER"
   policy_type = "DATABASE"
   re_auth = false
   source = "USER"
   mode = "COMMON"
   operator = "100032669045"
   id = 102606
 }
}

`
