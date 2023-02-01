package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverBackupMigrationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverBackupMigration,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_backup_migration.backup_migration", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_backup_migration.backup_migration",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverBackupMigration = `

resource "tencentcloud_sqlserver_backup_migration" "backup_migration" {
  instance_id = ""
  recovery_type = ""
  upload_type = ""
  migration_name = ""
  backup_files = 
}

`
