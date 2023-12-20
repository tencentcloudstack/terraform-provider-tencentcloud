package cfw_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCfwEdgeFirewallSwitchResource_basic -v
func TestAccTencentCloudNeedFixCfwEdgeFirewallSwitchResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwEdgeFirewallSwitch,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_firewall_switch.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_firewall_switch.example", "public_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_firewall_switch.example", "switch_mode"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_firewall_switch.example", "enable"),
				),
			},
			{
				Config: testAccCfwEdgeFirewallSwitchUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_firewall_switch.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_firewall_switch.example", "public_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_firewall_switch.example", "switch_mode"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_edge_firewall_switch.example", "enable"),
				),
			},
		},
	})
}

const testAccCfwEdgeFirewallSwitch = `
resource "tencentcloud_cfw_edge_firewall_switch" "example" {
  public_ip   = "43.138.44.22"
  switch_mode = 1
  enable      = 0
}
`

const testAccCfwEdgeFirewallSwitchUpdate = `
resource "tencentcloud_cfw_edge_firewall_switch" "example" {
  public_ip   = "43.138.44.22"
  switch_mode = 0
  enable      = 1
}
`
