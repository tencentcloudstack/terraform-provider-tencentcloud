package cfw_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCfwVpcFirewallSwitchResource_basic -v
func TestAccTencentCloudNeedFixCfwVpcFirewallSwitchResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwVpcFirewallSwitch,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_vpc_firewall_switch.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_vpc_firewall_switch.example", "vpc_ins_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_vpc_firewall_switch.example", "switch_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_vpc_firewall_switch.example", "enable"),
				),
			},
			{
				ResourceName:      "tencentcloud_cfw_vpc_firewall_switch.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCfwVpcFirewallSwitchUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_vpc_firewall_switch.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_vpc_firewall_switch.example", "vpc_ins_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_vpc_firewall_switch.example", "switch_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_vpc_firewall_switch.example", "enable"),
				),
			},
		},
	})
}

const testAccCfwVpcFirewallSwitch = `
resource "tencentcloud_cfw_vpc_firewall_switch" "example" {
  vpc_ins_id = "cfwg-c8c2de41"
  switch_id  = "cfws-f2c63ded84"
  enable     = 1
}
`

const testAccCfwVpcFirewallSwitchUpdate = `
resource "tencentcloud_cfw_vpc_firewall_switch" "example" {
  vpc_ins_id = "cfwg-c8c2de41"
  switch_id  = "cfws-f2c63ded84"
  enable     = 0
}
`
