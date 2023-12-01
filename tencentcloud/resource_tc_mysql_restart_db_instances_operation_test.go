package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlRestartDbInstancesOperationResource_basic -v
func TestAccTencentCloudMysqlRestartDbInstancesOperationResource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlRestartDbInstancesOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_restart_db_instances_operation.restart_db_instances_operation", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_restart_db_instances_operation.restart_db_instances_operation", "status", "1"),
				),
			},
		},
	})
}

const testAccMysqlRestartDbInstancesOperation = testAccMysqlInstanceEncryptionOperationVar + `

resource "tencentcloud_mysql_restart_db_instances_operation" "restart_db_instances_operation" {
  instance_id = tencentcloud_mysql_instance.mysql8.id
}

`
