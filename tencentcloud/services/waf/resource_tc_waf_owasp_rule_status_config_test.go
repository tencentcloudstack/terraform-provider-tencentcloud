package waf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWafOwaspRuleStatusConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafOwaspRuleStatusConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_owasp_rule_status_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_owasp_rule_status_config.example", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_owasp_rule_status_config.example", "rule_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_owasp_rule_status_config.example", "rule_status"),
				),
			},
			{
				Config: testAccWafOwaspRuleStatusConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_owasp_rule_status_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_owasp_rule_status_config.example", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_owasp_rule_status_config.example", "rule_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_owasp_rule_status_config.example", "rule_status"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_owasp_rule_status_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWafOwaspRuleStatusConfig = `
resource "tencentcloud_waf_owasp_rule_status_config" "example" {
  domain      = "example.qcloud.com"
  rule_id     = "106251141"
  rule_status = 0
}
`

const testAccWafOwaspRuleStatusConfigUpdate = `
resource "tencentcloud_waf_owasp_rule_status_config" "example" {
  domain      = "example.qcloud.com"
  rule_id     = "106251141"
  rule_status = 2
}
`
