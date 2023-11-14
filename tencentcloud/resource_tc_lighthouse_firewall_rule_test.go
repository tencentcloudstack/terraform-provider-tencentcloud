package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLighthouseFirewallRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseFirewallRule,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_firewall_rule.firewall_rule", "id")),
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
  instance_id = "lhins-acb1234"
  firewall_rules {
		protocol = "TCP"
		port = "80"
		cidr_block = "22"
		action = "ACCEPT"
		firewall_rule_description = "description"

  }
  firewall_version = 1
}

`
