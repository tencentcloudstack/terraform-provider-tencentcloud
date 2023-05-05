package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataSqlserverAccountsName = "data.tencentcloud_sqlserver_accounts.test"

func TestAccDataSourceTencentCloudSqlserverAccounts(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSqlserverAccountsBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataSqlserverAccountsName, "list.#"),
					resource.TestCheckResourceAttrSet(testDataSqlserverAccountsName, "list.0.status"),
					resource.TestCheckResourceAttrSet(testDataSqlserverAccountsName, "list.0.create_time"),
					resource.TestCheckResourceAttrSet(testDataSqlserverAccountsName, "list.0.update_time"),
					resource.TestCheckResourceAttrSet(testDataSqlserverAccountsName, "list.0.remark"),
				),
			},
		},
	})
}

const testAccTencentCloudDataSqlserverAccountsBasic = CommonPresetSQLServerAccount
