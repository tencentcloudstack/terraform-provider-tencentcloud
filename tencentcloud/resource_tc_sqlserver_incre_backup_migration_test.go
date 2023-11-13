package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverIncreBackupMigrationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverIncreBackupMigration,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_incre_backup_migration.incre_backup_migration", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_incre_backup_migration.incre_backup_migration",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverIncreBackupMigration = `

resource "tencentcloud_sqlserver_incre_backup_migration" "incre_backup_migration" {
  instance_id = "mssql-i1z41iwd"
  backup_migration_id = "migration_00001"
  backup_files = 
  is_recovery = "No"
}

`
