package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbRestartDBInstancesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbRestartDBInstances,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdb_restart_d_b_instances.restart_d_b_instances", "id")),
			},
			{
				ResourceName:      "tencentcloud_cdb_restart_d_b_instances.restart_d_b_instances",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdbRestartDBInstances = `

resource "tencentcloud_cdb_restart_d_b_instances" "restart_d_b_instances" {
  instance_ids = 
}

`
