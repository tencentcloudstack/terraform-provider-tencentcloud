package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbAccountPrivilegesResource_basic -v
func TestAccTencentCloudMariadbAccountPrivilegesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbAccountPrivileges,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mariadb_account_privileges.account_privileges", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_mariadb_account_privileges.account_privileges",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbAccountPrivileges = `
resource "tencentcloud_mariadb_account_privileges" "account_privileges" {
  instance_id = "tdsql-9vqvls95"
  accounts {
		user = "keep-modify-privileges"
		host = "127.0.0.1"
  }
  global_privileges = ["ALTER", "CREATE", "DELETE", "SELECT", "UPDATE", "DROP"]
}
`
