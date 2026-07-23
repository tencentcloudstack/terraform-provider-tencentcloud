package waf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWafApiSecSensitiveCustomEventRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafApiSecSensitiveCustomEventRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_api_sec_sensitive_custom_event_rule.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_waf_api_sec_sensitive_custom_event_rule.example", "rule_name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_waf_api_sec_sensitive_custom_event_rule.example", "status", "1"),
				),
			},
			{
				Config: testAccWafApiSecSensitiveCustomEventRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_waf_api_sec_sensitive_custom_event_rule.example", "status", "0"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_api_sec_sensitive_custom_event_rule.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWafApiSecSensitiveCustomEventRule = `
resource "tencentcloud_waf_api_sec_sensitive_custom_event_rule" "example" {
  domain        = "www.example.com"
  rule_name     = "tf-example"
  status        = 1
  description   = "tf example custom event rule"
  req_frequency = [100, 1]
  risk_level    = "300"
  source        = "custom"

  api_name_op {
    op    = "belong"
    value = ["/api/login"]

    api_name_method {
      api_name = "/api/login"
      method   = "POST"
    }
  }

  match_rule_list {
    key     = "get"
    operate = "contains"
    value   = ["admin"]
    name    = "role"
  }
}
`

const testAccWafApiSecSensitiveCustomEventRuleUpdate = `
resource "tencentcloud_waf_api_sec_sensitive_custom_event_rule" "example" {
  domain        = "www.example.com"
  rule_name     = "tf-example"
  status        = 0
  description   = "tf example custom event rule update"
  req_frequency = [200, 5]
  risk_level    = "200"
  source        = "custom"

  api_name_op {
    op    = "belong"
    value = ["/api/logout"]

    api_name_method {
      api_name = "/api/logout"
      method   = "GET"
    }
  }

  match_rule_list {
    key     = "post"
    operate = "equal"
    value   = ["guest"]
    name    = "role"
  }
}
`
