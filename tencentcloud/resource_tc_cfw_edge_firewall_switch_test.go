package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCfwEdgeFirewallSwitchResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwEdgeFirewallSwitch,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_firewall_switch.edge_firewall_switch", "id")),
			},
			{
				ResourceName:      "tencentcloud_cfw_edge_firewall_switch.edge_firewall_switch",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCfwEdgeFirewallSwitch = `

resource "tencentcloud_cfw_edge_firewall_switch" "edge_firewall_switch" {
  enable = 1
  edge_ip_switch_lst {
		public_ip = "1.1.1.1"
		subnet_id = "subnet-id"
		endpoint_ip = ""
		switch_mode = 0

  }
}

`
