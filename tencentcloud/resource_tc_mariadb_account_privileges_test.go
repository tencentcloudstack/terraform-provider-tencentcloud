package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

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
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mariadb_account_privileges.account_privileges", "id")),
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
  instance_id = "tdsql-e9tklsgz"
  accounts {
		user = ""
		host = ""

  }
  global_privileges = 
  database_privileges {
		privileges = 
		database = ""

  }
  table_privileges {
		database = ""
		table = ""
		privileges = 

  }
  column_privileges {
		database = ""
		table = ""
		column = ""
		privileges = 

  }
  view_privileges {
		database = ""
		view = ""
		privileges = 

  }
  function_privileges {
		database = ""
		function_name = ""
		privileges = 

  }
  procedure_privileges {
		database = ""
		procedure = ""
		privileges = 

  }
}

`
