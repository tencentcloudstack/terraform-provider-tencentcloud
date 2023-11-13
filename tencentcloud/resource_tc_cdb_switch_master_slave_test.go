package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbSwitchMasterSlaveResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbSwitchMasterSlave,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdb_switch_master_slave.switch_master_slave", "id")),
			},
			{
				ResourceName:      "tencentcloud_cdb_switch_master_slave.switch_master_slave",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdbSwitchMasterSlave = `

resource "tencentcloud_cdb_switch_master_slave" "switch_master_slave" {
  instance_id = ""
  dst_slave = ""
  force_switch = 
  wait_switch = 
}

`
