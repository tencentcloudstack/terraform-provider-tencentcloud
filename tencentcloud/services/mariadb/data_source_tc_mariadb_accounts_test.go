package mariadb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbAccountsDataSource -v
func TestAccTencentCloudMariadbAccountsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMariadbAccounts,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_accounts.accounts"),
				),
			},
		},
	})
}

const testAccDataSourceMariadbAccounts = testAccMariadbAccount + `

data "tencentcloud_mariadb_accounts" "accounts" {
	instance_id = tencentcloud_mariadb_hour_db_instance.basic.id
}

`
