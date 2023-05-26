package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverRestoreInstanceResource_basic -v
func TestAccTencentCloudSqlserverRestoreInstanceResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverRestoreInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_restore_instance.restore_instance", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_restore_instance.restore_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSqlserverConfigDeleteTmpDB,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_delete_db.config_delete_db", "id"),
				),
			},
		},
	})
}

const testAccSqlserverRestoreInstance = `
resource "tencentcloud_sqlserver_restore_instance" "restore_instance" {
  instance_id = "mssql-qelbzgwf"
  backup_id = 3461718019
  rename_restore {
    old_name = "keep_pubsub_db2"
  	new_name = "restore_keep_pubsub_db2"
  }
  type = 1
}
`

const testAccSqlserverConfigDeleteTmpDB = `
resource "tencentcloud_sqlserver_config_delete_db" "config_delete_db" {
  instance_id = "mssql-qelbzgwf"
  name = "restore_keep_pubsub_db2"
}
`
