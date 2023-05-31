package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMysqlSwitchMasterSlaveOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlSwitchMasterSlaveOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mysql_switch_master_slave_operation.switch_master_slave_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_mysql_switch_master_slave_operation.switch_master_slave_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMysqlSwitchMasterSlaveOperation = `

resource "tencentcloud_mysql_switch_master_slave_operation" "switch_master_slave_operation" {
  instance_id = ""
  dst_slave = ""
  force_switch = 
  wait_switch = 
}

`
