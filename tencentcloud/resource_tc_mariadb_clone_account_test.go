package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixMariadbCloneAccountResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbCloneAccount,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mariadb_clone_account.clone_account", "id"),
				),
			},
		},
	})
}

const testAccMariadbCloneAccount = `
resource "tencentcloud_mariadb_clone_account" "clone_account" {
  instance_id = "tdsql-9vqvls95"
  src_user = "srcuser"
  src_host = "10.13.1.%"
  dst_user = "dstuser"
  dst_host = "10.13.23.%"
  dst_desc = "test clone"
}
`
