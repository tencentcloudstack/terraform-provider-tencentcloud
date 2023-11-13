package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbAccountsResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbAccounts,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdb_accounts.accounts", "id")),
			},
			{
				ResourceName:      "tencentcloud_cdb_accounts.accounts",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdbAccounts = `

resource "tencentcloud_cdb_accounts" "accounts" {
  instance_id = ""
  accounts {
		user = ""
		host = ""

  }
  password = ""
  description = ""
  max_user_connections = 
}

`
