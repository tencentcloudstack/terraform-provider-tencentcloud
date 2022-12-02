package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func TestAccTencentCloudMysqlBackupPolicyResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTencentCloudMysqlBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlBackupPolicy(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudMysqlBackupPolicyExists("tencentcloud_mysql_backup_policy.mysql_backup_policy"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_backup_policy.mysql_backup_policy", "mysql_id"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_backup_policy.mysql_backup_policy", "retention_period", "56"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_backup_policy.mysql_backup_policy", "backup_time", "10:00-14:00"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_backup_policy.mysql_backup_policy", "binlog_period"),
				),
			},
			{
				Config: testAccMysqlBackupPolicyUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudMysqlBackupPolicyExists("tencentcloud_mysql_backup_policy.mysql_backup_policy"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_backup_policy.mysql_backup_policy", "mysql_id"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_backup_policy.mysql_backup_policy", "retention_period", "80"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_backup_policy.mysql_backup_policy", "backup_time", "06:00-10:00"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_backup_policy.mysql_backup_policy", "binlog_period"),
				),
			},
		},
	})
}

func testAccTencentCloudMysqlBackupPolicyExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		mysqlService := MysqlService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, err := mysqlService.DescribeBackupConfigByMysqlId(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccTencentCloudMysqlBackupPolicyDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	mysqlService := MysqlService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_mysql_backup_policy" {
			continue
		}

		instance, err := mysqlService.DescribeDBInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if instance == nil {
			return nil
		}

		desResponse, err := mysqlService.DescribeBackupConfigByMysqlId(ctx, rs.Primary.ID)

		if err != nil {
			sdkErr, ok := err.(*errors.TencentCloudSDKError)
			if ok && sdkErr.Code == MysqlInstanceIdNotFound {
				continue
			}
			return err
		}
		if *desResponse.Response.BackupExpireDays != 7 {
			return fmt.Errorf("mysql  backup  policy  BackupExpireDays  is not default")
		}
		if *desResponse.Response.BackupMethod != "physical" {
			return fmt.Errorf("mysql  backup  policy  BackupMethod  is not default")
		}
		if *desResponse.Response.StartTimeMax != 6 || *desResponse.Response.StartTimeMin != 2 {
			return fmt.Errorf("mysql  backup  policy  StartTimeMax or StartTimeMin  is not default")
		}
	}

	return nil
}

func testAccMysqlBackupPolicy() string {
	return fmt.Sprintf(`
%s
resource "tencentcloud_mysql_backup_policy" "mysql_backup_policy" {
  mysql_id         = local.mysql_id
  retention_period = 56
  backup_time      = "10:00-14:00"
}`, CommonPresetMysql)
}

func testAccMysqlBackupPolicyUpdate() string {
	return fmt.Sprintf(`
%s
resource "tencentcloud_mysql_backup_policy" "mysql_backup_policy" {
  mysql_id         = local.mysql_id
  retention_period = 80
  backup_time      = "06:00-10:00"
}`, CommonPresetMysql)
}
