package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverRollbackInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverRollbackInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_rollback_instance.rollback_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_rollback_instance.rollback_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverRollbackInstance = `

resource "tencentcloud_sqlserver_rollback_instance" "rollback_instance" {
  instance_id = "mssql-i1z41iwd"
  type = 0
  time = ""
  d_bs = 
  target_instance_id = ""
  rename_restore {
		old_name = ""
		new_name = ""

  }
}

`
