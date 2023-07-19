package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverDatabaseTDEResource_basic -v
func TestAccTencentCloudSqlserverDatabaseTDEResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverDatabaseTDE,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_database_tde.database_tde", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_database_tde.database_tde",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSqlserverDatabaseTDEUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_database_tde.database_tde", "id"),
				),
			},
		},
	})
}

const testAccSqlserverDatabaseTDE = `
resource "tencentcloud_sqlserver_database_tde" "database_tde" {
  instance_id = "mssql-qelbzgwf"
  db_names    = ["keep_tde_db", "keep_tde_db2"]
  encryption  = "enable"
}
`

const testAccSqlserverDatabaseTDEUpdate = `
resource "tencentcloud_sqlserver_database_tde" "database_tde" {
  instance_id = "mssql-qelbzgwf"
  db_names    = ["keep_tde_db", "keep_tde_db2"]
  encryption  = "disable"
}
`
