package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcUpdateRowFilterOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcUpdateRowFilterOperation,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_update_row_filter_operation.update_row_filter_operation", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_update_row_filter_operation.update_row_filter_operation", "policy_id", "103704"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_update_row_filter_operation.update_row_filter_operation", "policy.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_update_row_filter_operation.update_row_filter_operation", "policy_set.0.database"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_update_row_filter_operation.update_row_filter_operation", "policy_set.0.catalog"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_update_row_filter_operation.update_row_filter_operation", "policy_set.0.table", ""),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_update_row_filter_operation.update_row_filter_operation", "policy_set.0.operation"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_update_row_filter_operation.update_row_filter_operation", "policy_set.0.policy_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_update_row_filter_operation.update_row_filter_operation", "policy_set.0.re_auth"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_update_row_filter_operation.update_row_filter_operation", "policy_set.0.source"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_update_row_filter_operation.update_row_filter_operation", "policy_set.0.mode"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_update_row_filter_operation.update_row_filter_operation", "policy_set.0.id")),
			},
		},
	})
}

const testAccDlcUpdateRowFilterOperation = `

resource "tencentcloud_dlc_update_row_filter_operation" "update_row_filter_operation" {
  policy_id = 103704
  policy {
		database = "test_iac_keep"
		catalog = "DataLakeCatalog"
		table = "test_table"
		operation = "value!=\"0\""
		policy_type = "ROWFILTER"
		function = ""
		view = ""
		column = ""
		source = "USER"
		mode = "SENIOR"
        re_auth = false
  }
}

`
