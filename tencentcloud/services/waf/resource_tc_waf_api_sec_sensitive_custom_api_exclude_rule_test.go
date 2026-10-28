package waf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWafApiSecSensitiveCustomApiExcludeRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafApiSecSensitiveCustomApiExcludeRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_api_sec_sensitive_custom_api_exclude_rule.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_waf_api_sec_sensitive_custom_api_exclude_rule.example", "rule_name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_waf_api_sec_sensitive_custom_api_exclude_rule.example", "status", "1"),
				),
			},
			{
				Config: testAccWafApiSecSensitiveCustomApiExcludeRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_waf_api_sec_sensitive_custom_api_exclude_rule.example", "status", "0"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_api_sec_sensitive_custom_api_exclude_rule.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWafApiSecSensitiveCustomApiExcludeRule = `
resource "tencentcloud_waf_api_sec_sensitive_custom_api_exclude_rule" "example" {
  domain     = "www.example.com"
  rule_name  = "tf-example"
  status     = 1
  match_type = "prefix"
  content    = "/static/"
}
`

const testAccWafApiSecSensitiveCustomApiExcludeRuleUpdate = `
resource "tencentcloud_waf_api_sec_sensitive_custom_api_exclude_rule" "example" {
  domain     = "www.example.com"
  rule_name  = "tf-example"
  status     = 0
  match_type = "suffix"
  content    = ".js"
}
`
