package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverIncreBackupMigrationResource_basic -v
func TestAccTencentCloudSqlserverIncreBackupMigrationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverIncreBackupMigrationDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverIncreBackupMigration,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverIncreBackupMigrationExists("tencentcloud_sqlserver_incre_backup_migration.incre_backup_migration"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_incre_backup_migration.incre_backup_migration", "instance_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_incre_backup_migration.incre_backup_migration",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSqlserverIncreBackupMigrationUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverIncreBackupMigrationExists("tencentcloud_sqlserver_incre_backup_migration.incre_backup_migration"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_incre_backup_migration.incre_backup_migration", "instance_id"),
				),
			},
		},
	})
}

func testAccCheckSqlserverIncreBackupMigrationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_sqlserver_incre_backup_migration" {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken, id is %s", rs.Primary.ID)
		}

		instanceId := idSplit[0]
		backupMigrationId := idSplit[1]
		incrementalMigrationId := idSplit[2]

		result, err := service.DescribeSqlserverIncreBackupMigrationById(ctx, instanceId, backupMigrationId, incrementalMigrationId)
		if err != nil {
			return err
		}

		if result != nil {
			return fmt.Errorf("sqlserver incre backup migration %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckSqlserverIncreBackupMigrationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}

		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken, id is %s", rs.Primary.ID)
		}

		instanceId := idSplit[0]
		backupMigrationId := idSplit[1]
		incrementalMigrationId := idSplit[2]

		result, err := service.DescribeSqlserverIncreBackupMigrationById(ctx, instanceId, backupMigrationId, incrementalMigrationId)
		if err != nil {
			return err
		}

		if result == nil {
			return fmt.Errorf("sqlserver incre backup migration %s is not found", rs.Primary.ID)
		} else {
			return nil
		}
	}
}

const testAccSqlserverIncreBackupMigration = `
resource "tencentcloud_sqlserver_incre_backup_migration" "incre_backup_migration" {
  instance_id = "mssql-4gmc5805"
  backup_migration_id = "mssql-backup-migration-9tj0sxnz"
  backup_files = []
  is_recovery = "NO"
}
`

const testAccSqlserverIncreBackupMigrationUpdate = `
resource "tencentcloud_sqlserver_incre_backup_migration" "incre_backup_migration" {
  instance_id = "mssql-4gmc5805"
  backup_migration_id = "mssql-backup-migration-9tj0sxnz"
  backup_files = []
  is_recovery = "YES"
}
`
