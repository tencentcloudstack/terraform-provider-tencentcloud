package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverFullBackupMigrationResource_basic -v
func TestAccTencentCloudSqlserverFullBackupMigrationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverFullBackupMigrationDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverFullBackupMigration,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverFullBackupMigrationExists("tencentcloud_sqlserver_full_backup_migration.my_migration"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_full_backup_migration.my_migration", "instance_id"),
				),
			},
			{
				Config: testAccSqlserverFullBackupMigrationUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverFullBackupMigrationExists("tencentcloud_sqlserver_full_backup_migration.my_migration"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_full_backup_migration.my_migration", "instance_id"),
				),
			},
		},
	})
}

func testAccCheckSqlserverFullBackupMigrationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_sqlserver_full_backup_migration" {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		instanceId := rs.Primary.ID
		backupMigrationId := rs.Primary.Attributes["backup_migration_id"]

		result, err := service.DescribeSqlserverFullBackupMigrationById(ctx, instanceId, backupMigrationId)
		if err != nil {
			if sdkerr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "ResourceNotFound.InstanceNotFound" {
					return nil
				}
			}

			return err
		}

		if result != nil {
			return fmt.Errorf("sqlserver full_backup migration %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckSqlserverFullBackupMigrationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}

		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		instanceId := rs.Primary.ID
		backupMigrationId := rs.Primary.Attributes["backup_migration_id"]

		result, err := service.DescribeSqlserverFullBackupMigrationById(ctx, instanceId, backupMigrationId)
		if err != nil {
			return err
		}

		if result == nil {
			return fmt.Errorf("sqlserver full_backup migration %s is not found", rs.Primary.ID)
		} else {
			return nil
		}
	}
}

const testAccSqlserverFullBackupMigration = `
resource "tencentcloud_sqlserver_full_backup_migration" "my_migration" {
  instance_id = "mssql-qelbzgwf"
  recovery_type = "FULL"
  upload_type = "COS_URL"
  migration_name = "migration_test"
  backup_files = []
}
`

const testAccSqlserverFullBackupMigrationUpdate = `
resource "tencentcloud_sqlserver_full_backup_migration" "my_migration" {
  instance_id = "mssql-qelbzgwf"
  recovery_type = "FULL"
  upload_type = "COS_URL"
  migration_name = "migration_test_new"
  backup_files = []
}
`
