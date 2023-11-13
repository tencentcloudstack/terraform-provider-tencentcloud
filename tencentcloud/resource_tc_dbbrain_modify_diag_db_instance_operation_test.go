package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainModifyDiagDbInstanceOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainModifyDiagDbInstanceOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_modify_diag_db_instance_operation.modify_diag_db_instance_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_dbbrain_modify_diag_db_instance_operation.modify_diag_db_instance_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDbbrainModifyDiagDbInstanceOperation = `

resource "tencentcloud_dbbrain_modify_diag_db_instance_operation" "modify_diag_db_instance_operation" {
  instance_confs {
		daily_inspection = ""
		overview_display = ""

  }
  regions = ""
  product = ""
  instance_ids = 
}

`
