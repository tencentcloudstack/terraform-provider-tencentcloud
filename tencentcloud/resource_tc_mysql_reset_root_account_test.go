package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlResetRootAccountResource_basic -v
func TestAccTencentCloudMysqlResetRootAccountResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlResetRootAccount,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_reset_root_account.reset_root_account", "id"),
				),
			},
		},
	})
}

const testAccMysqlResetRootAccount = `

resource "tencentcloud_mysql_reset_root_account" "reset_root_account" {
	instance_id = "cdb-fitq5t9h"
}

`
