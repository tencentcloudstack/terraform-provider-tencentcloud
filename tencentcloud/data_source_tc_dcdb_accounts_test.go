package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDCDBAccountsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDcdbAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceDcdbAccounts, defaultDcdbInstanceId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_accounts.basic"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_accounts.basic", "list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_accounts.basic", "list.0.user_name"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_accounts.basic", "list.0.user_name", "mysql_ds"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_accounts.basic", "list.0.host", "127.0.0.1"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_accounts.basic", "list.0.description", "this is a test account"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_accounts.basic", "list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_accounts.basic", "list.0.update_time"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_accounts.basic", "list.0.read_only", "0"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_accounts.basic", "list.0.delay_thresh"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_accounts.basic", "list.0.slave_const"),
				),
			},
		},
	})
}

const testAccDataSourceDcdbAccounts = `
resource "tencentcloud_dcdb_account" "basic" {
	instance_id = "%s"
	user_name = "mysql_ds"
	host = "127.0.0.1"
	password = "===password==="
	read_only = 0
	description = "this is a test account"
	max_user_connections = 10
}

data "tencentcloud_dcdb_accounts" "basic" {
  instance_id = tencentcloud_dcdb_account.basic.instance_id
}

`
