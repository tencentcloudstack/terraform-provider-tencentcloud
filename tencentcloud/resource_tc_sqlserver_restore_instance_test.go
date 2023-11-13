package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverRestoreInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverRestoreInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_restore_instance.restore_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_restore_instance.restore_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverRestoreInstance = `

resource "tencentcloud_sqlserver_restore_instance" "restore_instance" {
  instance_id = "mssql-i1z41iwd"
  backup_id = 1981910
  target_instance_id = "mssql-au8ajamz"
  rename_restore {
		old_name = ""
		new_name = ""

  }
  type = 
  d_b_list = 
  group_id = ""
}

`
