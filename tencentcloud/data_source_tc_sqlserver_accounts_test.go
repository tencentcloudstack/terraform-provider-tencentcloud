package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataSqlserverAccountsName = "data.tencentcloud_sqlserver_accounts.test"

func TestAccDataSourceTencentCloudSqlserverAccounts(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSqlserverAccountsBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testDataSqlserverAccountsName, "list.#", "1"),
					resource.TestCheckResourceAttrSet(testDataSqlserverAccountsName, "list.0.status"),
					resource.TestCheckResourceAttrSet(testDataSqlserverAccountsName, "list.0.create_time"),
					resource.TestCheckResourceAttrSet(testDataSqlserverAccountsName, "list.0.update_time"),
					resource.TestCheckResourceAttrSet(testDataSqlserverAccountsName, "list.0.remark"),
				),
			},
		},
	})
}

const testAccTencentCloudDataSqlserverAccountsBasic = testAccSqlserverInstance + `
resource "tencentcloud_sqlserver_account" "test" {
  instance_id = tencentcloud_sqlserver_instance.test.id
  name = "tf_sqlserver_account1"
  password = "testt123"
}
data "tencentcloud_sqlserver_accounts" "test"{
  instance_id = tencentcloud_sqlserver_instance.test.id
  name = tencentcloud_sqlserver_account.test.name	
}
`
