package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	sdkError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func TestAccTencentCloudMysqlAccountResource(t *testing.T) {
	t.Parallel()
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
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		mysqlService := MysqlService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		split := strings.Split(rs.Primary.ID, FILED_SP)
		if len(split) < 2 {
			return fmt.Errorf("mysql account is not set")
		}

		var accountInfos []*cdb.AccountInfo
		var inErr, outErr error
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			accountInfos, inErr = mysqlService.DescribeAccounts(ctx, split[0])
			if inErr != nil {
				sdkErr, ok := inErr.(*sdkError.TencentCloudSDKError)
				if ok && sdkErr.Code == MysqlInstanceIdNotFound {
					return resource.NonRetryableError(fmt.Errorf("mysql account %s is not found", rs.Primary.ID))
				}
				return retryError(inErr, InternalError)

			}
			return nil
		})
		if outErr != nil {
			return outErr
		}
		for _, account := range accountInfos {
			if *account.User == split[1] && *account.Host == split[2] {
				return nil
			}
		}
		return fmt.Errorf("mysql account %s is not found", rs.Primary.ID)
	}
}

func testAccCheckMysqlAccountDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	mysqlService := MysqlService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, r := range s.RootModule().Resources {
		if r.Type != "tencentcloud_mysql_account" {
			continue
		}

		split := strings.Split(r.Primary.ID, FILED_SP)
		if len(split) < 2 {
			continue
		}
		instance, err := mysqlService.DescribeDBInstanceById(ctx, split[0])
		if err == nil && instance == nil {
			return nil
		}

		var accountInfos []*cdb.AccountInfo
		var inErr, outErr error
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			accountInfos, inErr = mysqlService.DescribeAccounts(ctx, split[0])
			if inErr != nil {
				sdkErr, ok := inErr.(*sdkError.TencentCloudSDKError)
				if ok && sdkErr.Code == MysqlInstanceIdNotFound {
					return nil
				}
				return retryError(inErr, InternalError)

			}
			return nil
		})

		if outErr != nil {
			return outErr
		}
		if accountInfos == nil {
			return nil
		}
		for _, account := range accountInfos {
			if *account.User == split[1] && *account.Host == split[2] {
				return fmt.Errorf("mysql account %s is still existing", split[1])
			}
		}
	}
	return nil
}

func testAccMysqlAccount() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_mysql_account" "mysql_account" {
	mysql_id = local.mysql_id
	name    = "test"
    host = "192.168.0.%%"
	password = "test1234"
	description = "test from terraform"
}
	`, CommonPresetMysql)
}
