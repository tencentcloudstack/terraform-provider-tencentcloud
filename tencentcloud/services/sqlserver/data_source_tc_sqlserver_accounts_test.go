package sqlserver_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataSqlserverAccountsName = "data.tencentcloud_sqlserver_accounts.test"

// go test -i; go test -test.run TestAccDataSourceTencentCloudSqlserverAccounts -v
func TestAccDataSourceTencentCloudSqlserverAccounts(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
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

const testAccTencentCloudDataSqlserverAccountsBasic = tcacctest.CommonPresetSQLServerAccount
