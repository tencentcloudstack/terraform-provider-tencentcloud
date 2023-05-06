package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func init() {
	resource.AddTestSweepers("tencentcloud_sqlserver_migration", &resource.Sweeper{
		Name: "tencentcloud_sqlserver_migration",
		F:    testSweepSqlserverMigration,
	})
}

// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_sqlserver_migration
func testSweepSqlserverMigration(r string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cli, _ := sharedClientForRegion(r)
	sqlServerService := SqlserverService{client: cli.(*TencentCloudClient).apiV3Conn}
	param := map[string]interface{}{}

	ret, err := sqlServerService.DescribeSqlserverMigrationsByFilter(ctx, param)
	if err != nil {
		return err
	}

	for _, v := range ret {
		delId := helper.UInt64ToStr(*v.MigrateId)

		if strings.HasPrefix(*v.MigrateName, keepResource) || strings.HasPrefix(*v.MigrateName, defaultResource) {
			continue
		}

		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			err := sqlServerService.DeleteSqlserverMigrationById(ctx, delId)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("[ERROR] sweeper tencentcloud_sqlserver_migration:[%v] failed! reason:[%s]", delId, err.Error())
		}
	}
	return nil
}

func TestAccTencentCloudSqlserverMigrationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSqlserverMigrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverMigration_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverMigrationExists("tencentcloud_sqlserver_migration.migration"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_migration.migration", "id"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_migration.migration", "migrate_name", "tf_test_migration"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_migration.migration", "migrate_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_migration.migration", "source_type", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_migration.migration", "migrate_db_set.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_migration.migration", "source.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_migration.migration", "target.#"),
				),
			},
			{
				Config: testAccSqlserverMigration_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverMigrationExists("tencentcloud_sqlserver_migration.migration"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_migration.migration", "id"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_migration.migration", "migrate_name", "tf_test_migration_changed"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_migration.migration", "migrate_type", "3"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_migration.migration",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSqlserverMigrationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_sqlserver_migration" {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		ret, err := service.DescribeSqlserverMigrationById(ctx, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("testAccCheckSqlserverMigrationDestroy found err:%v, id:%s", err.Error(), rs.Primary.ID)
		}
		if ret != nil {
			status := *ret.Status
			if status != SQLSERVER_MIGRATION_TERMINATED {
				return fmt.Errorf("SQL Server migration still exist, Id: %v, status:%v", rs.Primary.ID, status)
			}
		}
	}
	return nil
}

func testAccCheckSqlserverMigrationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		ret, err := service.DescribeSqlserverMigrationById(ctx, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("testAccCheckSqlserverMigrationExists found err:%v, id:%s", err.Error(), rs.Primary.ID)
		}
		if ret == nil {
			return fmt.Errorf("SQL Server migration not found, Id: %v", rs.Primary.ID)
		}
		return nil
	}
}

const testAccSqlserverMigration_account_db = testAccSqlserverBasicInstanceNetwork + CommonPresetSQLServerAccount + `
resource "tencentcloud_sqlserver_account" "src" {
	instance_id = local.sqlserver_id
	name = "tf_sqlserver_migration_src_account"
	password = "testt1234"
	is_admin = true
  }

  resource "tencentcloud_sqlserver_account_db_attachment" "src" {
	instance_id = local.sqlserver_id
	account_name = tencentcloud_sqlserver_account.src.name
	db_name = local.sqlserver_db # "keep_sqlserver_db"
	privilege = "ReadWrite"
  }

  resource "tencentcloud_sqlserver_instance" "dst" {
	name                          = "tf_sqlserver_dst_instance"
	availability_zone             = var.default_az
	charge_type                   = "POSTPAID_BY_HOUR"
	vpc_id                        = local.vpc_id
	subnet_id                     = local.subnet_id
	security_groups               = [local.sg_id]
	project_id                    = 0
	memory                        = 2
	storage                       = 10
	maintenance_week_set          = [1,2,3]
	maintenance_start_time        = "09:00"
	maintenance_time_span         = 3
	tags = {
	  "test"                      = "test"
	}
  }

  resource "tencentcloud_sqlserver_account" "dst" {
	instance_id = tencentcloud_sqlserver_instance.dst.id
	name = "tf_sqlserver_migration_dst_account"
	password = "testt1234"
	is_admin = true
  }

  resource "tencentcloud_sqlserver_db" "dst" {
	instance_id = tencentcloud_sqlserver_instance.dst.id
	name        = "tf_migration_dst_db"
	charset     = "Chinese_PRC_BIN"
	remark      = "testACC-remark"
  }
`

const testAccSqlserverMigration_basic = testAccSqlserverMigration_account_db + `

resource "tencentcloud_sqlserver_migration" "migration" {
  migrate_name = "tf_test_migration"
  migrate_type = 1
  source_type = 1
  source {
	instance_id = local.sqlserver_id
	user_name = tencentcloud_sqlserver_account.src.name
	password = tencentcloud_sqlserver_account.src.password
  }
  target {
	instance_id = tencentcloud_sqlserver_instance.dst.id
	user_name = tencentcloud_sqlserver_account.dst.name
	password = tencentcloud_sqlserver_account.dst.password
  }

  migrate_db_set {
	db_name = local.sqlserver_db
  }
}
`

const testAccSqlserverMigration_update = testAccSqlserverMigration_account_db + `

resource "tencentcloud_sqlserver_migration" "migration" {
  migrate_name = "tf_test_migration_changed"
  migrate_type = 3
  source_type = 1
  source {
	instance_id = local.sqlserver_id
	user_name = tencentcloud_sqlserver_account.src.name
	password = tencentcloud_sqlserver_account.src.password
  }
  target {
	instance_id = tencentcloud_sqlserver_instance.dst.id
	user_name = tencentcloud_sqlserver_account.dst.name
	password = tencentcloud_sqlserver_account.dst.password
  }

  migrate_db_set {
	db_name = local.sqlserver_db
  }
}
`
