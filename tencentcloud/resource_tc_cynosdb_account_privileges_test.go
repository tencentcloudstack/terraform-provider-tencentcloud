package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbAccountPrivilegesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbAccountPrivileges,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_account_privileges.account_privileges", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_account_privileges.account_privileges",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbAccountPrivileges = `

resource "tencentcloud_cynosdb_account_privileges" "account_privileges" {
  cluster_id = "xxx"
  account {
		account_name = ""
		host = ""

  }
  global_privileges = 
  database_privileges {
		db = ""
		privileges = 

  }
  table_privileges {
		db = ""
		table_name = ""
		privileges = 

  }
}

`
