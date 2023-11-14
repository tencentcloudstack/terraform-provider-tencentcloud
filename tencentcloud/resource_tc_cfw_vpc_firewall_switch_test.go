package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCfwVpcFirewallSwitchResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwVpcFirewallSwitch,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cfw_vpc_firewall_switch.vpc_firewall_switch", "id")),
			},
			{
				ResourceName:      "tencentcloud_cfw_vpc_firewall_switch.vpc_firewall_switch",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCfwVpcFirewallSwitch = `

resource "tencentcloud_cfw_vpc_firewall_switch" "vpc_firewall_switch" {
  enable = 1
  all_switch = 0
  switch_list {
		switch_mode = 1
		switch_id = "cfws-id"

  }
}

`
