package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverBackupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverBackup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_backup.backup", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_backup.backup",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverBackup = `

resource "tencentcloud_sqlserver_backup" "backup" {
  strategy = 
  db_names = 
  instance_id = ""
  backup_name = ""
  backup_id = ""
  group_id = ""
}

`
