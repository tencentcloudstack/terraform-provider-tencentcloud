package waf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWafApiSecSensitiveCustomRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafApiSecSensitiveCustomRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_api_sec_sensitive_custom_rule.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_waf_api_sec_sensitive_custom_rule.example", "rule_name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_waf_api_sec_sensitive_custom_rule.example", "status", "1"),
				),
			},
			{
				Config: testAccWafApiSecSensitiveCustomRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_waf_api_sec_sensitive_custom_rule.example", "status", "0"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_api_sec_sensitive_custom_rule.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWafApiSecSensitiveCustomRule = `
resource "tencentcloud_waf_api_sec_sensitive_custom_rule" "example" {
  domain      = "www.example.com"
  rule_name   = "tf-example"
  status      = 1
  position    = ["query"]
  match_key   = "keyword"
  match_value = ["password"]
  level       = "300"
  match_cond  = ["contains"]
  is_pan      = 0
}
`

const testAccWafApiSecSensitiveCustomRuleUpdate = `
resource "tencentcloud_waf_api_sec_sensitive_custom_rule" "example" {
  domain      = "www.example.com"
  rule_name   = "tf-example"
  status      = 0
  position    = ["query"]
  match_key   = "keyword"
  match_value = ["password", "secret"]
  level       = "200"
  match_cond  = ["contains"]
  is_pan      = 1
}
`
