package cfw_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCfwNatFirewallSwitchResource_basic -v
func TestAccTencentCloudNeedFixCfwNatFirewallSwitchResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwNatFirewallSwitch,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_firewall_switch.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_firewall_switch.example", "nat_ins_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_firewall_switch.example", "subnet_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_firewall_switch.example", "enable"),
				),
			},
			{
				ResourceName:      "tencentcloud_cfw_nat_firewall_switch.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCfwNatFirewallSwitchUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_firewall_switch.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_firewall_switch.example", "nat_ins_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_firewall_switch.example", "subnet_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_firewall_switch.example", "enable"),
				),
			},
		},
	})
}

const testAccCfwNatFirewallSwitch = `
resource "tencentcloud_cfw_nat_firewall_switch" "example" {
  nat_ins_id = "cfwnat-18d2ba18"
  subnet_id  = "subnet-ef7wyymr"
  enable     = 1
}
`

const testAccCfwNatFirewallSwitchUpdate = `
resource "tencentcloud_cfw_nat_firewall_switch" "example" {
  nat_ins_id = "cfwnat-18d2ba18"
  subnet_id  = "subnet-ef7wyymr"
  enable     = 0
}
`
