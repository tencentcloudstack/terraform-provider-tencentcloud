package waf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWafApiSecSensitiveSceneRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafApiSecSensitiveSceneRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_api_sec_sensitive_scene_rule.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_waf_api_sec_sensitive_scene_rule.example", "rule_name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_waf_api_sec_sensitive_scene_rule.example", "status", "1"),
				),
			},
			{
				Config: testAccWafApiSecSensitiveSceneRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_waf_api_sec_sensitive_scene_rule.example", "status", "0"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_api_sec_sensitive_scene_rule.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWafApiSecSensitiveSceneRule = `
resource "tencentcloud_waf_api_sec_sensitive_scene_rule" "example" {
  domain    = "www.example.com"
  rule_name = "tf-example"
  status    = 1
  source    = "custom"

  rule_list {
    key     = "get"
    operate = "contains"
    value   = ["login"]
    name    = "username"
  }
}
`

const testAccWafApiSecSensitiveSceneRuleUpdate = `
resource "tencentcloud_waf_api_sec_sensitive_scene_rule" "example" {
  domain    = "www.example.com"
  rule_name = "tf-example"
  status    = 0
  source    = "custom"

  rule_list {
    key     = "post"
    operate = "equal"
    value   = ["register"]
    name    = "account"
  }
}
`
