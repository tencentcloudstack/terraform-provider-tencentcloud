package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverConfigDeleteDBResource_basic -v
func TestAccTencentCloudSqlserverConfigDeleteDBResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverDBDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverTmpInstanceDB,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_db.test", "id"),
				),
			},
			{
				Config: testAccSqlserverConfigDeleteDB,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_delete_db.config_delete_db", "id"),
				),
			},
		},
	})
}

const testAccSqlserverTmpInstanceDB = `
resource "tencentcloud_sqlserver_db" "test" {
  instance_id = "mssql-qelbzgwf"
  name        = "create_db_test"
  charset     = "Chinese_PRC_BIN"
  remark      = "test-remark"
}
`

const testAccSqlserverConfigDeleteDB = `
resource "tencentcloud_sqlserver_config_delete_db" "config_delete_db" {
  instance_id = "mssql-qelbzgwf"
  name = "create_db_test"
}
`
