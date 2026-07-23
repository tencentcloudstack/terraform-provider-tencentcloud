package waf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWafApiSecSensitiveWhiteRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafApiSecSensitiveWhiteRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_api_sec_sensitive_white_rule.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_waf_api_sec_sensitive_white_rule.example", "rule_name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_waf_api_sec_sensitive_white_rule.example", "status", "1"),
				),
			},
			{
				Config: testAccWafApiSecSensitiveWhiteRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_waf_api_sec_sensitive_white_rule.example", "status", "0"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_api_sec_sensitive_white_rule.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWafApiSecSensitiveWhiteRule = `
resource "tencentcloud_waf_api_sec_sensitive_white_rule" "example" {
  domain      = "www.example.com"
  rule_name   = "tf-example"
  status      = 1
  white_mode  = 2
  description = "tf example white rule"

  api_name_op {
    op    = "belong"
    value = ["/api/user/info"]

    api_name_method {
      api_name = "/api/user/info"
      method   = "GET"
    }
  }

  white_fields {
    field_name      = "id_card"
    field_type      = "body"
    sensitive_types = ["IDCARD"]
  }
}
`

const testAccWafApiSecSensitiveWhiteRuleUpdate = `
resource "tencentcloud_waf_api_sec_sensitive_white_rule" "example" {
  domain      = "www.example.com"
  rule_name   = "tf-example"
  status      = 0
  white_mode  = 1
  description = "tf example white rule update"

  api_name_op {
    op    = "belong"
    value = ["/api/user/profile"]

    api_name_method {
      api_name = "/api/user/profile"
      method   = "POST"
    }
  }
}
`
