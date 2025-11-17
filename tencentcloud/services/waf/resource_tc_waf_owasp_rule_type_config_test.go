package waf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWafOwaspRuleTypeConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafOwaspRuleTypeConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_owasp_rule_type_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_owasp_rule_type_config.example", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_owasp_rule_type_config.example", "type_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_owasp_rule_type_config.example", "rule_type_status"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_owasp_rule_type_config.example", "rule_type_action"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_owasp_rule_type_config.example", "rule_type_level"),
				),
			},
			{
				Config: testAccWafOwaspRuleTypeConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_owasp_rule_type_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_owasp_rule_type_config.example", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_owasp_rule_type_config.example", "type_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_owasp_rule_type_config.example", "rule_type_status"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_owasp_rule_type_config.example", "rule_type_action"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_owasp_rule_type_config.example", "rule_type_level"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_owasp_rule_type_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWafOwaspRuleTypeConfig = `
resource "tencentcloud_waf_owasp_rule_type_config" "example" {
  domain           = "example.qcloud.com"
  type_id          = "30000000"
  rule_type_status = 1
  rule_type_action = 1
  rule_type_level  = 200
}
`

const testAccWafOwaspRuleTypeConfigUpdate = `
resource "tencentcloud_waf_owasp_rule_type_config" "example" {
  domain           = "example.qcloud.com"
  type_id          = "30000000"
  rule_type_status = 0
  rule_type_action = 0
  rule_type_level  = 100
}
`
