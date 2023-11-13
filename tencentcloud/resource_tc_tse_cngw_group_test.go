package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTseCngwGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseCngwGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_group.cngw_group", "id")),
			},
			{
				ResourceName:      "tencentcloud_tse_cngw_group.cngw_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTseCngwGroup = `

resource "tencentcloud_tse_cngw_group" "cngw_group" {
  gateway_id = ""
  name = ""
  node_config {
		specification = ""
		number = 

  }
  subnet_id = ""
  description = ""
  internet_max_bandwidth_out = 
  internet_config {
		internet_address_version = ""
		internet_pay_mode = ""
		internet_max_bandwidth_out = 
		description = ""
		sla_type = ""
		multi_zone_flag = 
		master_zone_id = ""
		slave_zone_id = ""

  }
}

`
