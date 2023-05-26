package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverRollbackInstanceResource_basic -v
func TestAccTencentCloudSqlserverRollbackInstanceResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverRollbackInstance(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_rollback_instance.rollback_instance", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_rollback_instance.rollback_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSqlserverConfigDeleteRollBackDB,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_delete_db.config_delete_db", "id"),
				),
			},
		},
	})
}

const testAccSqlserverRollbackInstanceString = `
resource "tencentcloud_sqlserver_rollback_instance" "rollback_instance" {
  instance_id = "mssql-qelbzgwf"
  type = 1
  time = "%s"
  dbs = ["keep_pubsub_db2"]
  rename_restore {
    old_name = "keep_pubsub_db2"
	new_name = "rollback_pubsub_db3"
  }
}`

func testAccSqlserverRollbackInstance() string {
	var currentTime = time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04:05")
	return fmt.Sprintf(testAccSqlserverRollbackInstanceString, currentTime)
}

const testAccSqlserverConfigDeleteRollBackDB = `
resource "tencentcloud_sqlserver_config_delete_db" "config_delete_db" {
  instance_id = "mssql-qelbzgwf"
  name = "rollback_pubsub_db3"
}
`
