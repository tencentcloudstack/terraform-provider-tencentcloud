package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataSqlserverAccountsName = "data.tencentcloud_sqlserver_accounts.test"

func TestAccTencentCloudDataSqlserverAccounts(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSqlserverAccountsBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataSqlserverAccountsName, "list.0.status"),
					resource.TestCheckResourceAttrSet(testDataSqlserverAccountsName, "list.0.create_time"),
					resource.TestCheckResourceAttrSet(testDataSqlserverAccountsName, "list.0.update_time"),
					resource.TestCheckResourceAttrSet(testDataSqlserverAccountsName, "list.0.remark"),
				),
			},
		},
	})
}

const testAccTencentCloudDataSqlserverAccountsBasic = testAccSqlserverDB_basic + `
data "tencentcloud_sqlserver_accounts" "test"{
	  instance_id = tencentcloud_sqlserver_instance.test.id
}
`
