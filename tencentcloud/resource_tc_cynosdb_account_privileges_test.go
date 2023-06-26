package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbAccountPrivilegesResource_basic -v
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_account_privileges.account_privileges", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_account_privileges.account_privileges", "account_name"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account_privileges.account_privileges", "global_privileges.#", "5"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account_privileges.account_privileges", "database_privileges.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account_privileges.account_privileges", "database_privileges.0.db", "users"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account_privileges.account_privileges", "database_privileges.0.privileges.#", "6"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account_privileges.account_privileges", "table_privileges.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account_privileges.account_privileges", "table_privileges.0.db", "users"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account_privileges.account_privileges", "table_privileges.0.table_name", "tb_user_name"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account_privileges.account_privileges", "table_privileges.0.privileges.#", "3"),
				),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_account_privileges.account_privileges",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCynosdbAccountPrivilegesUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_account_privileges.account_privileges", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_account_privileges.account_privileges", "account_name"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account_privileges.account_privileges", "global_privileges.#", "4"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account_privileges.account_privileges", "database_privileges.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account_privileges.account_privileges", "database_privileges.0.db", "users"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account_privileges.account_privileges", "database_privileges.0.privileges.#", "5"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account_privileges.account_privileges", "table_privileges.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account_privileges.account_privileges", "table_privileges.0.db", "users"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account_privileges.account_privileges", "table_privileges.0.table_name", "tb_user_name"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account_privileges.account_privileges", "table_privileges.0.privileges.#", "2"),
				),
			},
		},
	})
}

const testAccCynosdbAccountPrivilegesVar = CommonCynosdb + `

resource "tencentcloud_cynosdb_account" "account" {
	cluster_id = var.cynosdb_cluster_id
	account_name = "terraform_account"
	account_password = "Password@1234"
	host = "%"
	description = "terraform test"
	max_user_connections = 2
}

`

const testAccCynosdbAccountPrivileges = testAccCynosdbAccountPrivilegesVar + `

resource "tencentcloud_cynosdb_account_privileges" "account_privileges" {
	cluster_id   = var.cynosdb_cluster_id
	account_name = tencentcloud_cynosdb_account.account.account_name
	host         = "%"
	global_privileges = [
	  "CREATE",
	  "DROP",
	  "ALTER",
	  "CREATE TEMPORARY TABLES",
	  "CREATE VIEW"
	]
	database_privileges {
	  db = "users"
	  privileges = [
		"DROP",
		"REFERENCES",
		"INDEX",
		"CREATE VIEW",
		"INSERT",
		"EVENT"
	  ]
	}
	table_privileges {
	  db         = "users"
	  table_name = "tb_user_name"
	  privileges = [
		"ALTER",
		"REFERENCES",
		"SHOW VIEW"
	  ]
	}
}

`

const testAccCynosdbAccountPrivilegesUp = testAccCynosdbAccountPrivilegesVar + `

resource "tencentcloud_cynosdb_account_privileges" "account_privileges" {
	cluster_id   = var.cynosdb_cluster_id
	account_name = tencentcloud_cynosdb_account.account.account_name
	host         = "%"
	global_privileges = [
	  "CREATE",
	  "DROP",
	  "ALTER",
	  "CREATE VIEW"
	]
	database_privileges {
	  db = "users"
	  privileges = [
		"DROP",
		"REFERENCES",
		"INDEX",
		"CREATE VIEW",
		"EVENT"
	  ]
	}
	table_privileges {
	  db         = "users"
	  table_name = "tb_user_name"
	  privileges = [
		"ALTER",
		"SHOW VIEW"
	  ]
	}
}

`
