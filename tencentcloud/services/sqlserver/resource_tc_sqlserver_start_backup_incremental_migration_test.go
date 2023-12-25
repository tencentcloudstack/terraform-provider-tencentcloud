package sqlserver_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixSqlserverStartBackupIncrementalMigrationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverStartBackupIncrementalMigration,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_start_backup_incremental_migration.start_backup_incremental_migration", "id"),
				),
			},
		},
	})
}

const testAccSqlserverStartBackupIncrementalMigration = `
resource "tencentcloud_sqlserver_start_backup_incremental_migration" "start_backup_incremental_migration" {
  instance_id              = "mssql-i1z41iwd"
  backup_migration_id      = "mssql-backup-migration-cg0ffgqt"
  incremental_migration_id = "mssql-incremental-migration-kp7bgv8p"
}
`
