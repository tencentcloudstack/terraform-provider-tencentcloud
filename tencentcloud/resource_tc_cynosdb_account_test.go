package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbAccountResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbAccount,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_account.account", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_account.account",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbAccount = `

resource "tencentcloud_cynosdb_account" "account" {
  cluster_id = "xxx"
  accounts {
		account_name = ""
		account_password = ""
		host = ""
		description = ""
		max_user_connections = 

  }
}

`
