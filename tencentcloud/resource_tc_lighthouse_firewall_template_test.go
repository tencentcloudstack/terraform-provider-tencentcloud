package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseFirewallTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseFirewallTemplate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_firewall_template.firewall_template", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_lighthouse_firewall_template.firewall_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLighthouseFirewallTemplate = `

resource "tencentcloud_lighthouse_firewall_template" "firewall_template" {
	template_name = "firewall-template-test"
	template_rules {
		protocol = "TCP"
		port = "8080"
		cidr_block = "127.0.0.1"
		action = "ACCEPT"
		firewall_rule_description = "test description"
	}
	template_rules {
		protocol = "TCP"
		port = "8090"
		cidr_block = "127.0.0.0/24"
		action = "DROP"
		firewall_rule_description = "test description"
	}
}
`
