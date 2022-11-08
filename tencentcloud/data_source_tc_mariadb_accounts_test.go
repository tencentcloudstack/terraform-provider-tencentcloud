package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbAccountsDataSource -v
func TestAccTencentCloudMariadbAccountsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMariadbAccounts,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_accounts.accounts"),
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
