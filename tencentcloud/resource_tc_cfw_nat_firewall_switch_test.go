package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCfwNatFirewallSwitchResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwNatFirewallSwitch,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_firewall_switch.nat_firewall_switch", "id")),
			},
			{
				ResourceName:      "tencentcloud_cfw_nat_firewall_switch.nat_firewall_switch",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCfwNatFirewallSwitch = `

resource "tencentcloud_cfw_nat_firewall_switch" "nat_firewall_switch" {
  enable = 1
  cfw_ins_id_list = 
  subnet_id_list = 
  route_table_id_list = 
}

`
