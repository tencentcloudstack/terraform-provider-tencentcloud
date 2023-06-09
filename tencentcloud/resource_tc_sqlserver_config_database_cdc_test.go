package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverConfigDatabaseCDCResource_basic -v
func TestAccTencentCloudSqlserverConfigDatabaseCDCResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverConfigDatabaseCDC,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_database_cdc.config_database_cdc", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_config_database_cdc.config_database_cdc",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverConfigDatabaseCDC = `
resource "tencentcloud_sqlserver_config_database_cdc" "config_database_cdc" {
  db_name = "keep_pubsub_db2"
  modify_type = "disable"
  instance_id = "mssql-qelbzgwf"
}
`
