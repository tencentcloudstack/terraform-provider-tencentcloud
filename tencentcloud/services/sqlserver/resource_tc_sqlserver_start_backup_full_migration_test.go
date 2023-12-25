package sqlserver_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixSqlserverStartBackupFullMigrationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverStartBackupFullMigration,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_start_backup_full_migration.start_backup_full_migration", "id"),
				),
			},
		},
	})
}

const testAccSqlserverStartBackupFullMigration = `
resource "tencentcloud_sqlserver_start_backup_full_migration" "start_backup_full_migration" {
  instance_id         = "mssql-i1z41iwd"
  backup_migration_id = "mssql-backup-migration-kpl74n9l"
}
`
