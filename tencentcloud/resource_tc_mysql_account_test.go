package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func TestAccTencentCloudMysqlAccountResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMysqlAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlAccount(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlAccountExists("tencentcloud_mysql_account.mysql_account"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_account.mysql_account", "mysql_id"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_account.mysql_account", "name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_account.mysql_account", "description", "test from terraform"),
				),
			},
		},
	})
}

func testAccCheckMysqlAccountExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := GetLogId(nil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		mysqlService := MysqlService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		split := strings.Split(rs.Primary.ID, FILED_SP)
		if len(split) < 2 {
			return fmt.Errorf("mysql account is not set")
		}
		accounts, err := mysqlService.DescribeAccounts(ctx, split[0])
		if err != nil {
			return err
		}
		for _, account := range accounts {
			if *account.User == split[1] {
				return nil
			}
		}
		return fmt.Errorf("mysql account %s is not found", rs.Primary.ID)
	}
}

func testAccCheckMysqlAccountDestroy(s *terraform.State) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	mysqlService := MysqlService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, r := range s.RootModule().Resources {
		if r.Type != "tencentcloud_mysql_account" {
			continue
		}

		split := strings.Split(r.Primary.ID, FILED_SP)
		if len(split) < 2 {
			continue
		}
		accounts, err := mysqlService.DescribeAccounts(ctx, split[0])
		if err != nil {
			sdkErr, ok := err.(*errors.TencentCloudSDKError)
			if ok && sdkErr.Code == MysqlInstanceIdNotFound {
				continue
			}
			return err
		}
		for _, account := range accounts {
			if *account.User == split[1] {
				return fmt.Errorf("mysql account %s is still existing", split[1])
			}
		}
	}
	return nil
}

func testAccMysqlAccount() string {
	return fmt.Sprintf(`
resource "tencentcloud_mysql_instance" "mysql" {
	pay_type = 1
	mem_size = 1000
	volume_size = 50
	instance_name = "testAccMysqlAccount"
	vpc_id = "vpc-fzdzrsir"
	subnet_id = "subnet-he8ldxx6"
	engine_version = "5.7"
	root_password = "test1234"
	availability_zone = "ap-guangzhou-4"
}

resource "tencentcloud_mysql_account" "mysql_account" {
	mysql_id = "${tencentcloud_mysql_instance.mysql.id}"
	name = "test"
	password = "test1234"
	description = "test from terraform"
}
	`)
}
