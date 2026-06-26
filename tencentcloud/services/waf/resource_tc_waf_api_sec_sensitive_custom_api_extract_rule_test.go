package waf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWafApiSecSensitiveCustomApiExtractRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafApiSecSensitiveCustomApiExtractRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_api_sec_sensitive_custom_api_extract_rule.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_waf_api_sec_sensitive_custom_api_extract_rule.example", "rule_name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_waf_api_sec_sensitive_custom_api_extract_rule.example", "status", "1"),
				),
			},
			{
				Config: testAccWafApiSecSensitiveCustomApiExtractRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_waf_api_sec_sensitive_custom_api_extract_rule.example", "status", "0"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_api_sec_sensitive_custom_api_extract_rule.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWafApiSecSensitiveCustomApiExtractRule = `
resource "tencentcloud_waf_api_sec_sensitive_custom_api_extract_rule" "example" {
  domain    = "www.example.com"
  rule_name = "tf-example"
  status    = 1
  api_name  = "/api/login"
  methods   = ["GET", "POST"]
  regex     = "/api/.*"
}
`

const testAccWafApiSecSensitiveCustomApiExtractRuleUpdate = `
resource "tencentcloud_waf_api_sec_sensitive_custom_api_extract_rule" "example" {
  domain    = "www.example.com"
  rule_name = "tf-example"
  status    = 0
  api_name  = "/api/logout"
  methods   = ["GET"]
  regex     = "/api/v2/.*"
}
`
