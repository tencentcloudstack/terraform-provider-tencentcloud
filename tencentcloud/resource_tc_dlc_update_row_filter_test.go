package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcUpdateRowFilterResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcUpdateRowFilter,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_update_row_filter.update_row_filter", "id")),
			},
		},
	})
}

const testAccDlcUpdateRowFilter = `

resource "tencentcloud_dlc_update_row_filter" "update_row_filter" {
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
