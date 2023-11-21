package tencentcloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

// go test -i; go test -test.run TestAccTencentCloudTestingMysqlAccountResource_basic -v
func TestAccTencentCloudTestingMysqlAccountResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMysqlAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTestingMysqlAccount(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlAccountExists("tencentcloud_mysql_account.mysql_account"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_account.mysql_account", "mysql_id"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_account.mysql_account", "name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_account.mysql_account", "description", "test from terraform"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_account.mysql_account", "max_user_connections", "10"),
				),
			},
			{
				Config: testAccTestingMysqlAccountUp(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlAccountExists("tencentcloud_mysql_account.mysql_account"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_account.mysql_account", "mysql_id"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_account.mysql_account", "name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_account.mysql_account", "description", "test from terraform"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_account.mysql_account", "max_user_connections", "10"),
				),
			},
		},
	})
}

func testAccTestingMysqlAccount() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_mysql_account" "mysql_account" {
	mysql_id = local.mysql_id
	name    = "terraform_test"
    host = "192.168.0.%%"
	password = "Test@123456#"
	description = "test from terraform"
	max_user_connections = 10
}
	`, CommonPresetMysql)
}

func testAccTestingMysqlAccountUp() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_mysql_account" "mysql_account" {
	mysql_id = local.mysql_id
	name    = "terraform_test"
    host = "192.168.1.%%"
	password = "Test@123456#"
	description = "test from terraform"
	max_user_connections = 10
}
	`, CommonPresetMysql)
}
