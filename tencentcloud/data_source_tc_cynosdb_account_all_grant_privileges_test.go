package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbAccountAllGrantPrivilegesDataSource_basic -v
func TestAccTencentCloudCynosdbAccountAllGrantPrivilegesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbAccountAllGrantPrivilegesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_account_all_grant_privileges.account_all_grant_privileges"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_account_all_grant_privileges.account_all_grant_privileges", "database_privileges.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_account_all_grant_privileges.account_all_grant_privileges", "global_privileges.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_account_all_grant_privileges.account_all_grant_privileges", "privilege_statements.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_account_all_grant_privileges.account_all_grant_privileges", "table_privileges.#"),
				),
			},
		},
	})
}

const testAccCynosdbAccountAllGrantPrivilegesDataSource = `
data "tencentcloud_cynosdb_account_all_grant_privileges" "account_all_grant_privileges" {
  cluster_id = "cynosdbmysql-bws8h88b"
  account {
    account_name = "keep_dts"
    host         = "%"
  }
}
`
