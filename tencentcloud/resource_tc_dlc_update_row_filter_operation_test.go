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
					resource.TestCheckResourceAttr("tencentcloud_dlc_update_row_filter_operation.update_row_filter_operation", "policy.0.database", "test_iac_keep"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_update_row_filter_operation.update_row_filter_operation", "policy.0.catalog", "DataLakeCatalog"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_update_row_filter_operation.update_row_filter_operation", "policy.0.table", "test_table"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_update_row_filter_operation.update_row_filter_operation", "policy.0.operation", "value!=\"0\""),
					resource.TestCheckResourceAttr("tencentcloud_dlc_update_row_filter_operation.update_row_filter_operation", "policy.0.policy_type", "ROWFILTER"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_update_row_filter_operation.update_row_filter_operation", "policy.0.re_auth", "false"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_update_row_filter_operation.update_row_filter_operation", "policy.0.source", "USER"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_update_row_filter_operation.update_row_filter_operation", "policy.0.mode", "SENIOR")),
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
