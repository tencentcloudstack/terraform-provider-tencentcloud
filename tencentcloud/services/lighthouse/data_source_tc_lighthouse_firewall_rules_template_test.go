package lighthouse_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseFirewallRulesTemplateDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseFirewallRulesTemplateDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_lighthouse_firewall_rules_template.firewall_rules_template")),
			},
		},
	})
}

const testAccLighthouseFirewallRulesTemplateDataSource = `

data "tencentcloud_lighthouse_firewall_rules_template" "firewall_rules_template" {
  }

`
