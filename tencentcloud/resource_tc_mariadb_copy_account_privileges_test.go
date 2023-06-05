package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbCopyAccountPrivilegesResource_basic -v
func TestAccTencentCloudMariadbCopyAccountPrivilegesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbCopyAccountPrivileges,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mariadb_copy_account_privileges.copy_account_privileges", "id"),
				),
			},
		},
	})
}

const testAccMariadbCopyAccountPrivileges = `
resource "tencentcloud_mariadb_copy_account_privileges" "copy_account_privileges" {
  instance_id   = "tdsql-9vqvls95"
  src_user_name = "keep-modify-privileges"
  src_host      = "127.0.0.1"
  dst_user_name = "keep-copy-user"
  dst_host      = "127.0.0.1"
}
`
