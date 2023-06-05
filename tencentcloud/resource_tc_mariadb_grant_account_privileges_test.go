package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbGrantAccountPrivilegesResource_basic -v
func TestAccTencentCloudMariadbGrantAccountPrivilegesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbGrantAccountPrivileges,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mariadb_grant_account_privileges.grant_account_privileges", "id"),
				),
			},
		},
	})
}

const testAccMariadbGrantAccountPrivileges = `
resource "tencentcloud_mariadb_grant_account_privileges" "grant_account_privileges" {
  instance_id = "tdsql-9vqvls95"
  user_name   = "keep-grant"
  host        = "127.0.0.1"
  db_name     = "*"
  privileges  = ["SELECT", "INSERT"]
}
`
