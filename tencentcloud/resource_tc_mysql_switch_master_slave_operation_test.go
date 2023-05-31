package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlSwitchMasterSlaveOperationResource_basic -v
func TestAccTencentCloudMysqlSwitchMasterSlaveOperationResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlSwitchMasterSlaveOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_switch_master_slave_operation.switch_master_slave_operation", "id"),
				),
			},
		},
	})
}

const testAccMysqlSwitchMasterSlaveOperationVar = `
variable "instance_id" {
  default = "` + defaultDbBrainInstanceId + `"
}
`

const testAccMysqlSwitchMasterSlaveOperation = testAccMysqlSwitchMasterSlaveOperationVar + `

resource "tencentcloud_mysql_switch_master_slave_operation" "switch_master_slave_operation" {
	instance_id = var.instance_id
	dst_slave = "first"
	force_switch = true
	wait_switch = true
}

`
