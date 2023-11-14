package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverFullBackupMigrationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverFullBackupMigration,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_full_backup_migration.full_backup_migration", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_full_backup_migration.full_backup_migration",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverFullBackupMigration = `

resource "tencentcloud_sqlserver_full_backup_migration" "full_backup_migration" {
  instance_id = "mssql-i1z41iwd"
  recovery_type = "FULL"
  upload_type = "COS_URL"
  migration_name = "test_migration"
  backup_files = 
}

`
