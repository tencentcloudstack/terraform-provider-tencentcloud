package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverConfigDatabaseMdfResource_basic -v
func TestAccTencentCloudSqlserverConfigDatabaseMdfResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverConfigDatabaseMdf,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_database_mdf.config_database_mdf", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_config_database_mdf.config_database_mdf",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverConfigDatabaseMdf = `
resource "tencentcloud_sqlserver_config_database_mdf" "config_database_mdf" {
  db_name = "keep_pubsub_db2"
  instance_id = "mssql-qelbzgwf"
}
`
