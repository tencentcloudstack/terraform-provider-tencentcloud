package sqlserver_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcsqlserver "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/sqlserver"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverFullBackupMigrationResource_basic -v
func TestAccTencentCloudSqlserverFullBackupMigrationResource_basic(t *testing.T) {
	t.Parallel()
	loc, _ := time.LoadLocation("Asia/Chongqing")
	startTime := time.Now().AddDate(0, 0, -3).In(loc).Format("2006-01-02 15:04:05")
	endTime := time.Now().AddDate(0, 0, 1).In(loc).Format("2006-01-02 15:04:05")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverFullBackupMigrationDestroy,
		Providers:    tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccSqlserverFullBackupMigration, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverFullBackupMigrationExists("tencentcloud_sqlserver_full_backup_migration.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_full_backup_migration.example", "instance_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_full_backup_migration.my_migration",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: fmt.Sprintf(testAccSqlserverFullBackupMigrationUpdate, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverFullBackupMigrationExists("tencentcloud_sqlserver_full_backup_migration.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_full_backup_migration.example", "instance_id"),
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := svcsqlserver.NewSqlserverService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken, id is %s", rs.Primary.ID)
		}

		instanceId := idSplit[0]
		backupMigrationId := idSplit[1]

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

		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := svcsqlserver.NewSqlserverService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken, id is %s", rs.Primary.ID)
		}

		instanceId := idSplit[0]
		backupMigrationId := idSplit[1]

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

const testAccSqlserverFullBackupMigration = tcacctest.DefaultVpcSubnets + tcacctest.DefaultSecurityGroupData + `
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

data "tencentcloud_sqlserver_backups" "example" {
  instance_id = tencentcloud_sqlserver_db.example.instance_id
  backup_name = tencentcloud_sqlserver_general_backup.example.backup_name
  start_time  = "%s"
  end_time    = "%s"
}

resource "tencentcloud_sqlserver_basic_instance" "example" {
  name                   = "tf-example"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = local.vpc_id
  subnet_id              = local.subnet_id
  project_id             = 0
  memory                 = 4
  storage                = 100
  cpu                    = 2
  machine_type           = "CLOUD_PREMIUM"
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  security_groups        = [local.sg_id]

  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_sqlserver_db" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
  name        = "tf_example_db"
  charset     = "Chinese_PRC_BIN"
  remark      = "test-remark"
}

resource "tencentcloud_sqlserver_general_backup" "example" {
  instance_id = tencentcloud_sqlserver_db.example.instance_id
  backup_name = "tf_example_backup"
  strategy    = 0
}

resource "tencentcloud_sqlserver_full_backup_migration" "example" {
  instance_id    = tencentcloud_sqlserver_basic_instance.example.id
  recovery_type  = "FULL"
  upload_type    = "COS_URL"
  migration_name = "migration_test"
  backup_files   = [data.tencentcloud_sqlserver_backups.example.list.0.internet_url]
}
`

const testAccSqlserverFullBackupMigrationUpdate = tcacctest.DefaultVpcSubnets + tcacctest.DefaultSecurityGroupData + `
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

data "tencentcloud_sqlserver_backups" "example" {
  instance_id = tencentcloud_sqlserver_db.example.instance_id
  backup_name = tencentcloud_sqlserver_general_backup.example.backup_name
  start_time  = "%s"
  end_time    = "%s"
}

resource "tencentcloud_sqlserver_basic_instance" "example" {
  name                   = "tf-example"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = local.vpc_id
  subnet_id              = local.subnet_id
  project_id             = 0
  memory                 = 4
  storage                = 100
  cpu                    = 2
  machine_type           = "CLOUD_PREMIUM"
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  security_groups        = [local.sg_id]

  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_sqlserver_db" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
  name        = "tf_example_db"
  charset     = "Chinese_PRC_BIN"
  remark      = "test-remark"
}

resource "tencentcloud_sqlserver_general_backup" "example" {
  instance_id = tencentcloud_sqlserver_db.example.instance_id
  backup_name = "tf_example_backup"
  strategy    = 0
}

resource "tencentcloud_sqlserver_full_backup_migration" "example" {
  instance_id    = tencentcloud_sqlserver_basic_instance.example.id
  recovery_type  = "FULL"
  upload_type    = "COS_URL"
  migration_name = "migration_test_update"
  backup_files   = [data.tencentcloud_sqlserver_backups.example.list.0.internet_url]
}
`
