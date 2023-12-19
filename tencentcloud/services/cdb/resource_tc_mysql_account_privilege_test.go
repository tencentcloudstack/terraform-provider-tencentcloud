package cdb_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	sdkError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	localcdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cdb"
)

func TestAccTencentCloudMysqlAccountPrivilegeResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccMysqlAccountPrivilegeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlAccountPrivilege(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccMysqlAccountPrivilegeExists("tencentcloud_mysql_account_privilege.mysql_account_privilege"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_account_privilege.mysql_account_privilege", "mysql_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_account_privilege.mysql_account_privilege", "account_name"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_account_privilege.mysql_account_privilege", "database_names.#", "1"),

					resource.TestCheckResourceAttr("tencentcloud_mysql_account_privilege.mysql_account_privilege", "privileges.#", "4"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mysql_account_privilege.mysql_account_privilege", "privileges.*", "SELECT"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mysql_account_privilege.mysql_account_privilege", "privileges.*", "UPDATE"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mysql_account_privilege.mysql_account_privilege", "privileges.*", "INSERT"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mysql_account_privilege.mysql_account_privilege", "privileges.*", "DELETE"),
				),
			},
			{
				Config: testAccMysqlAccountPrivilegeUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccMysqlAccountPrivilegeExists("tencentcloud_mysql_account_privilege.mysql_account_privilege"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_account_privilege.mysql_account_privilege", "mysql_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_account_privilege.mysql_account_privilege", "account_name"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_account_privilege.mysql_account_privilege", "database_names.#", "1"),

					resource.TestCheckResourceAttr("tencentcloud_mysql_account_privilege.mysql_account_privilege", "privileges.#", "1"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mysql_account_privilege.mysql_account_privilege", "privileges.*", "TRIGGER"),
				),
			},
		},
	})
}

func testAccMysqlAccountPrivilegeExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		mysqlService := localcdb.NewMysqlService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		var privilegeId localcdb.ResourceTencentCloudMysqlAccountPrivilegeId

		if err := json.Unmarshal([]byte(rs.Primary.ID), &privilegeId); err != nil {
			return fmt.Errorf("Local data[terraform.tfstate] corruption,can not got old account privilege id")
		}

		var inErr, outErr error

		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, inErr = mysqlService.DescribeAccountPrivileges(ctx, privilegeId.MysqlId, privilegeId.AccountName, privilegeId.AccountHost, []string{"test"})
			if inErr != nil {
				if sdkErr, ok := inErr.(*sdkError.TencentCloudSDKError); ok {
					if sdkErr.Code == localcdb.MysqlInstanceIdNotFound {
						return resource.NonRetryableError(fmt.Errorf("privilege not exists in mysql"))
					}
					if sdkErr.Code == "InvalidParameter" && strings.Contains(sdkErr.GetMessage(), "instance not found") {
						return resource.NonRetryableError(fmt.Errorf("privilege not exists in mysql"))
					}
					if sdkErr.Code == "InternalError.TaskError" && strings.Contains(sdkErr.Message, "User does not exist") {
						return resource.NonRetryableError(fmt.Errorf("privilege not exists in mysql"))
					}

				}
			}
			return nil
		})

		if outErr != nil {
			return outErr
		}

		var accountInfos []*cdb.AccountInfo
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			accountInfos, inErr = mysqlService.DescribeAccounts(ctx, privilegeId.MysqlId)
			if inErr != nil {
				sdkErr, ok := inErr.(*sdkError.TencentCloudSDKError)
				if ok && sdkErr.Code == localcdb.MysqlInstanceIdNotFound {
					return resource.NonRetryableError(fmt.Errorf("mysql account %s is not found", rs.Primary.ID))
				}
				return tccommon.RetryError(inErr, tccommon.InternalError)

			}
			return nil
		})
		if outErr != nil {
			return outErr
		}
		for _, account := range accountInfos {
			if *account.User == privilegeId.AccountName && *account.Host == privilegeId.AccountHost {
				return nil
			}
		}
		return fmt.Errorf("mysql  aacount privilege not found on server")
	}

}

func testAccMysqlAccountPrivilegeDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	mysqlService := localcdb.NewMysqlService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_mysql_account_privilege" {
			continue
		}
		var privilegeId localcdb.ResourceTencentCloudMysqlAccountPrivilegeId

		if err := json.Unmarshal([]byte(rs.Primary.ID), &privilegeId); err != nil {
			return fmt.Errorf("Local data[terraform.tfstate] corruption,can not got old account privilege id")
		}

		instance, err := mysqlService.DescribeDBInstanceById(ctx, privilegeId.MysqlId)
		if err == nil && instance == nil {
			return nil
		}

		var privileges []string
		var inErr, outErr error

		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			privileges, inErr = mysqlService.DescribeAccountPrivileges(ctx, privilegeId.MysqlId, privilegeId.AccountName, privilegeId.AccountHost, []string{"test"})
			if inErr != nil {
				if sdkErr, ok := inErr.(*sdkError.TencentCloudSDKError); ok {
					if sdkErr.Code == localcdb.MysqlInstanceIdNotFound {
						return nil
					}
					if sdkErr.Code == "InvalidParameter" && strings.Contains(sdkErr.GetMessage(), "instance not found") {
						return nil
					}
					if sdkErr.Code == "InternalError.TaskError" && strings.Contains(sdkErr.Message, "User does not exist") {
						return nil
					}

				}
			}
			return nil
		})

		if outErr != nil {
			return outErr
		}

		if len(privileges) == 0 {
			return nil
		}

		if len(privileges) != 1 || privileges[0] != localcdb.MYSQL_DATABASE_MUST_PRIVILEGE {
			return fmt.Errorf("mysql  aacount privilege not clean ok")
		}
	}

	return nil
}

func testAccMysqlAccountPrivilege() string {
	return fmt.Sprintf(`
%s
resource "tencentcloud_mysql_account" "mysql_account" {
  mysql_id    = local.mysql_id
  name        = "previlege_test"
  host        = "119.168.110.%%"
  password    = "test1234"
  description = "test from terraform"
}
resource "tencentcloud_mysql_account_privilege" "mysql_account_privilege" {
  mysql_id       = local.mysql_id
  account_name   = tencentcloud_mysql_account.mysql_account.name
  account_host   = tencentcloud_mysql_account.mysql_account.host
  privileges     = ["SELECT", "INSERT", "UPDATE", "DELETE"]
  database_names = ["test"]
}`, tcacctest.CommonPresetMysql)
}

func testAccMysqlAccountPrivilegeUpdate() string {
	return fmt.Sprintf(`
%s
resource "tencentcloud_mysql_account" "mysql_account" {
  mysql_id    = local.mysql_id
  name        = "previlege_test"
  host        = "119.168.110.%%"
  password    = "test1234"
  description = "test from terraform"
}
resource "tencentcloud_mysql_account_privilege" "mysql_account_privilege" {
  mysql_id       = local.mysql_id
  account_name   = tencentcloud_mysql_account.mysql_account.name
  account_host   = tencentcloud_mysql_account.mysql_account.host
  privileges     = ["TRIGGER"]
  database_names = ["test"]
}`, tcacctest.CommonPresetMysql)

}
