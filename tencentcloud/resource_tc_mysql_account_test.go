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
				Config: testAccMysqlAccount(MysqlInstanceCommonTestCase),
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
		logId := getLogId(nil)
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
	logId := getLogId(nil)
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

func testAccMysqlAccount(commonTestCase string) string {
	return fmt.Sprintf(`
%s
resource "tencentcloud_mysql_account" "mysql_account" {
	mysql_id = "${tencentcloud_mysql_instance.default.id}"
	name = "test"
	password = "test1234"
	description = "test from terraform"
}
	`, commonTestCase)
}
