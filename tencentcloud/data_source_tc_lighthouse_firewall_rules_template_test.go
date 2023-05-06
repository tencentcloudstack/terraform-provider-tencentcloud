package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseFirewallRulesTemplateDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseFirewallRulesTemplateDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_lighthouse_firewall_rules_template.firewall_rules_template")),
			},
		},
	})
}

const testAccLighthouseFirewallRulesTemplateDataSource = `

data "tencentcloud_lighthouse_firewall_rules_template" "firewall_rules_template" {
  }

`
