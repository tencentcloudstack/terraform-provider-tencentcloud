package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainModifyDiagDbInstanceConfResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainModifyDiagDbInstanceConf,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_modify_diag_db_instance_conf.modify_diag_db_instance_conf", "id")),
			},
			{
				ResourceName:      "tencentcloud_dbbrain_modify_diag_db_instance_conf.modify_diag_db_instance_conf",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDbbrainModifyDiagDbInstanceConf = `

resource "tencentcloud_dbbrain_modify_diag_db_instance_conf" "modify_diag_db_instance_conf" {
  instance_confs {
		daily_inspection = ""
		overview_display = ""

  }
  regions = ""
  product = ""
  instance_ids = 
}

`
