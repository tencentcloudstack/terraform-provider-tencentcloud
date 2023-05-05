package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseFirewallRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseFirewallRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_firewall_rule.firewall_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_lighthouse_firewall_rule.firewall_rule", "firewall_rules.0.cidr_block", "10.0.0.1"),
					resource.TestCheckResourceAttr("tencentcloud_lighthouse_firewall_rule.firewall_rule", "firewall_rules.1.cidr_block", "10.0.0.2"),
				),
			},
			{
				Config: testAccLighthouseFirewallRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_firewall_rule.firewall_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_lighthouse_firewall_rule.firewall_rule", "firewall_rules.0.cidr_block", "10.0.0.1"),
					resource.TestCheckResourceAttr("tencentcloud_lighthouse_firewall_rule.firewall_rule", "firewall_rules.1.cidr_block", "10.0.0.3"),
				),
			},
			{
				ResourceName:      "tencentcloud_lighthouse_firewall_rule.firewall_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLighthouseFirewallRule = `

resource "tencentcloud_lighthouse_firewall_rule" "firewall_rule" {
  instance_id = "lhins-hwe21u91"
  firewall_rules {
	protocol = "TCP"
	port = "80"
	cidr_block = "10.0.0.1"
	action = "ACCEPT"
	firewall_rule_description = "description 1"
  }
  firewall_rules {
	protocol = "TCP"
	port = "80"
	cidr_block = "10.0.0.2"
	action = "ACCEPT"
	firewall_rule_description = "description 2"
  }
}
`

const testAccLighthouseFirewallRuleUpdate = `

resource "tencentcloud_lighthouse_firewall_rule" "firewall_rule" {
  instance_id = "lhins-hwe21u91"
  firewall_rules {
	protocol = "TCP"
	port = "80"
	cidr_block = "10.0.0.1"
	action = "ACCEPT"
	firewall_rule_description = "description 1"
  }
  firewall_rules {
	protocol = "TCP"
	port = "80"
	cidr_block = "10.0.0.3"
	action = "ACCEPT"
	firewall_rule_description = "description 2"
  }
}
`
