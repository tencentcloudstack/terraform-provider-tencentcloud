package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverStartBackupIncrementalMigrationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverStartBackupIncrementalMigration,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_start_backup_incremental_migration.start_backup_incremental_migration", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_start_backup_incremental_migration.start_backup_incremental_migration",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverStartBackupIncrementalMigration = `

resource "tencentcloud_sqlserver_start_backup_incremental_migration" "start_backup_incremental_migration" {
  instance_id = "mssql-i1z41iwd"
  backup_migration_id = ""
  incremental_migration_id = ""
}

`
