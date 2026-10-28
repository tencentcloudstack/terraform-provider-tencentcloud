package waf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWafApiSecSensitivePrivilegeRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafApiSecSensitivePrivilegeRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_api_sec_sensitive_privilege_rule.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_waf_api_sec_sensitive_privilege_rule.example", "rule_name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_waf_api_sec_sensitive_privilege_rule.example", "status", "1"),
				),
			},
			{
				Config: testAccWafApiSecSensitivePrivilegeRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_waf_api_sec_sensitive_privilege_rule.example", "status", "0"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_api_sec_sensitive_privilege_rule.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWafApiSecSensitivePrivilegeRule = `
resource "tencentcloud_waf_api_sec_sensitive_privilege_rule" "example" {
  domain         = "www.example.com"
  rule_name      = "tf-example"
  status         = 1
  api_name       = ["/api/user/info"]
  position       = "header"
  parameter_list = ["token"]
  option         = 1

  api_name_op {
    op    = "belong"
    value = ["/api/user/info"]

    api_name_method {
      api_name = "/api/user/info"
      method   = "GET"
    }
  }
}
`

const testAccWafApiSecSensitivePrivilegeRuleUpdate = `
resource "tencentcloud_waf_api_sec_sensitive_privilege_rule" "example" {
  domain         = "www.example.com"
  rule_name      = "tf-example"
  status         = 0
  api_name       = ["/api/user/profile"]
  position       = "cookie"
  parameter_list = ["session"]
  option         = 1

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
