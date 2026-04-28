package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoSecurityJSInjectionRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoSecurityJSInjectionRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_security_js_injection_rule.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_js_injection_rule.example", "js_injection_rules.0.name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_js_injection_rule.example", "js_injection_rules.0.priority", "50"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_js_injection_rule.example", "js_injection_rules.0.inject_js", "inject-sdk-only"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_security_js_injection_rule.example", "js_injection_rules.0.rule_id"),
				),
			},
			{
				Config: testAccTeoSecurityJSInjectionRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_teo_security_js_injection_rule.example", "js_injection_rules.0.name", "tf-example-update"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_js_injection_rule.example", "js_injection_rules.0.priority", "60"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_security_js_injection_rule.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoSecurityJSInjectionRule = `
resource "tencentcloud_teo_security_js_injection_rule" "example" {
  zone_id = "zone-3fkff38fyw8s"

  js_injection_rules {
    name      = "tf-example"
    priority  = 50
    condition = "$${http.request.host} in ['test.makn.cn']"
    inject_js = "inject-sdk-only"
  }
}
`

const testAccTeoSecurityJSInjectionRuleUpdate = `
resource "tencentcloud_teo_security_js_injection_rule" "example" {
  zone_id = "zone-3fkff38fyw8s"

  js_injection_rules {
    name      = "tf-example-update"
    priority  = 60
    condition = "$${http.request.host} in ['test.makn.cn']"
    inject_js = "inject-sdk-only"
  }
}
`
