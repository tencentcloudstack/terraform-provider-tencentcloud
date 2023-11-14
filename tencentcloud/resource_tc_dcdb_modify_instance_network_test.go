package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcdbModifyInstanceNetworkResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbModifyInstanceNetwork,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dcdb_modify_instance_network.modify_instance_network", "id")),
			},
			{
				ResourceName:      "tencentcloud_dcdb_modify_instance_network.modify_instance_network",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcdbModifyInstanceNetwork = `

resource "tencentcloud_dcdb_modify_instance_network" "modify_instance_network" {
  instance_id = ""
  vpc_id = ""
  subnet_id = ""
  vip = ""
  vipv6 = ""
  vip_release_delay = 
}

`
