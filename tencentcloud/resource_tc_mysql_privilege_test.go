package tencentcloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/hashcode"

	"strconv"
	"strings"

	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

var testAccTencentCloudMysqlPrivilegeType = "tencentcloud_mysql_privilege"
var testAccTencentCloudMysqlPrivilegeName = testAccTencentCloudMysqlPrivilegeType + ".privilege"

func TestAccTencentCloudMysqlPrivilege(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMysqlPrivilegeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlPrivilege(mysqlInstanceCommonTestCase),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccMysqlPrivilegeExists,
					resource.TestCheckResourceAttrSet(testAccTencentCloudMysqlPrivilegeName, "mysql_id"),
					resource.TestCheckResourceAttrSet(testAccTencentCloudMysqlPrivilegeName, "account_name"),
					resource.TestCheckResourceAttr(testAccTencentCloudMysqlPrivilegeName, "global.#", "1"),
					resource.TestCheckResourceAttr(testAccTencentCloudMysqlPrivilegeName, "table.#", "1"),
					resource.TestCheckResourceAttr(testAccTencentCloudMysqlPrivilegeName, "column.#", "1"),
					resource.TestCheckResourceAttr(testAccTencentCloudMysqlPrivilegeName, "global."+strconv.Itoa(hashcode.String("TRIGGER")), "TRIGGER"),
				),
			},
			{
				Config: testAccMysqlPrivilegeUpdate(mysqlInstanceCommonTestCase),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccMysqlPrivilegeExists,
					resource.TestCheckResourceAttrSet(testAccTencentCloudMysqlPrivilegeName, "mysql_id"),
					resource.TestCheckResourceAttrSet(testAccTencentCloudMysqlPrivilegeName, "account_name"),
					resource.TestCheckResourceAttr(testAccTencentCloudMysqlPrivilegeName, "global."+strconv.Itoa(hashcode.String("TRIGGER")), "TRIGGER"),

					//diff
					resource.TestCheckResourceAttr(testAccTencentCloudMysqlPrivilegeName, "global.#", "2"),
					resource.TestCheckResourceAttr(testAccTencentCloudMysqlPrivilegeName, "table.#", "2"),
					resource.TestCheckResourceAttr(testAccTencentCloudMysqlPrivilegeName, "column.#", "0"),
					resource.TestCheckResourceAttr(testAccTencentCloudMysqlPrivilegeName, "global."+strconv.Itoa(hashcode.String("SELECT")), "SELECT"),
				),
			},
		},
	})
}

func testAccMysqlPrivilegeExists(s *terraform.State) error {

	rs, ok := s.RootModule().Resources[testAccTencentCloudMysqlPrivilegeName]
	if !ok {
		return fmt.Errorf("resource %s is not found", testAccTencentCloudMysqlPrivilegeName)
	}

	var privilegeId resourceTencentCloudMysqlPrivilegeId

	if err := json.Unmarshal([]byte(rs.Primary.ID), &privilegeId); err != nil {
		return fmt.Errorf("Local data[terraform.tfstate] corruption,can not got old account privilege id")
	}

	request := cdb.NewDescribeAccountPrivilegesRequest()
	request.InstanceId = &privilegeId.MysqlId
	request.User = &privilegeId.AccountName
	request.Host = stringToPointer(MYSQL_DEFAULT_ACCOUNT_HOST)

	response, err := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn.UseMysqlClient().DescribeAccountPrivileges(request)
	if err != nil {
		return err
	}
	if response == nil || response.Response == nil {
		return errors.New("sdk DescribeAccountPrivileges return error,miss Response")
	}

	if len(response.Response.GlobalPrivileges) > 0 ||
		len(response.Response.ColumnPrivileges) > 0 ||
		len(response.Response.TablePrivileges) > 0 ||
		len(response.Response.DatabasePrivileges) > 0 {
		return nil
	}
	return fmt.Errorf("set privilege return nil")
}

func testAccMysqlPrivilegeDestroy(s *terraform.State) error {
	rs, ok := s.RootModule().Resources[testAccTencentCloudMysqlPrivilegeName]
	if !ok {
		return fmt.Errorf("resource %s is not found", testAccTencentCloudMysqlPrivilegeName)
	}

	var privilegeId resourceTencentCloudMysqlPrivilegeId

	if err := json.Unmarshal([]byte(rs.Primary.ID), &privilegeId); err != nil {
		return fmt.Errorf("Local data[terraform.tfstate] corruption,can not got old account privilege id")
	}

	mysqlService := MysqlService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	instance, err := mysqlService.DescribeDBInstanceById(contextNil, privilegeId.MysqlId)
	if instance == nil {
		return nil
	}

	request := cdb.NewDescribeAccountPrivilegesRequest()
	request.InstanceId = &privilegeId.MysqlId
	request.User = &privilegeId.AccountName
	request.Host = stringToPointer(MYSQL_DEFAULT_ACCOUNT_HOST)

	response, err := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn.UseMysqlClient().DescribeAccountPrivileges(request)
	if err != nil {
		if sdkerr, ok := err.(*sdkErrors.TencentCloudSDKError); ok && sdkerr.GetCode() == "InvalidParameter" &&
			strings.Contains(sdkerr.GetMessage(), "instance not found") {
			return nil
		}
		return err
	}

	if response == nil || response.Response == nil {
		return errors.New("sdk DescribeAccountPrivileges return error,miss Response")
	}

	if len(response.Response.GlobalPrivileges) > 0 ||
		len(response.Response.ColumnPrivileges) > 0 ||
		len(response.Response.TablePrivileges) > 0 ||
		len(response.Response.DatabasePrivileges) > 0 {
		return fmt.Errorf("privilege is still exist")
	}
	return nil
}

func testAccMysqlPrivilege(commonTestCase string) string {
	return fmt.Sprintf(`
%s
resource "tencentcloud_mysql_account" "mysql_account" {
  mysql_id    = tencentcloud_mysql_instance.default.id
  name        = "test11"
  password    = "test1234"
  description = "test from terraform"
}

resource "tencentcloud_mysql_privilege" "privilege" {
  mysql_id     = tencentcloud_mysql_instance.default.id
  account_name = tencentcloud_mysql_account.mysql_account.name
  global       = ["TRIGGER"]
  database {
    privileges    = ["SELECT"]
    database_name = "performance_schema"
  }
  table {
    privileges    = ["SELECT", "INSERT", "UPDATE"]
    database_name = "mysql"
    table_name    = "user"
  }
  column {
    privileges    = ["SELECT"]
    database_name = "mysql"
    table_name    = "user"
    column_name   = "host"
  }
}`, commonTestCase)
}

func testAccMysqlPrivilegeUpdate(commonTestCase string) string {
	return fmt.Sprintf(`
%s
resource "tencentcloud_mysql_account" "mysql_account" {
  mysql_id    = tencentcloud_mysql_instance.default.id
  name        = "test11"
  password    = "test1234"
  description = "test from terraform"
}

resource "tencentcloud_mysql_privilege" "privilege" {
  mysql_id     = tencentcloud_mysql_instance.default.id
  account_name = tencentcloud_mysql_account.mysql_account.name
  global       = ["TRIGGER","SELECT"]
  table {
    privileges    = ["SELECT"]
    database_name = "mysql"
    table_name    = "user"
  }
  table {
    privileges    = ["SELECT"]
    database_name = "mysql"
    table_name    = "db"
  }
}`, commonTestCase)
}
