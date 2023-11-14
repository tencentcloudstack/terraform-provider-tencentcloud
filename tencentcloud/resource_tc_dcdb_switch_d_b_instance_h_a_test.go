package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcdbSwitchDBInstanceHAResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbSwitchDBInstanceHA,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dcdb_switch_d_b_instance_h_a.switch_d_b_instance_h_a", "id")),
			},
			{
				ResourceName:      "tencentcloud_dcdb_switch_d_b_instance_h_a.switch_d_b_instance_h_a",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcdbSwitchDBInstanceHA = `

resource "tencentcloud_dcdb_switch_d_b_instance_h_a" "switch_d_b_instance_h_a" {
  instance_id = ""
  zone = ""
}

`
