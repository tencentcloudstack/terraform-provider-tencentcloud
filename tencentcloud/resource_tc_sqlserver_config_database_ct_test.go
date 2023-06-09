package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverConfigDatabaseCTResource_basic -v
func TestAccTencentCloudSqlserverConfigDatabaseCTResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverConfigDatabaseCT,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_database_ct.config_database_ct", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_config_database_ct.config_database_ct",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverConfigDatabaseCT = `
resource "tencentcloud_sqlserver_config_database_ct" "config_database_ct" {
  db_name = "keep_pubsub_db2"
  modify_type = "enable"
  instance_id = "mssql-qelbzgwf"
  change_retention_day = 7
}
`
