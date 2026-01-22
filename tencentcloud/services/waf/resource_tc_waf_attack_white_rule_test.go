package waf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWafAttackWhiteRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWafAttackWhiteRule,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_waf_attack_white_rule.waf_attack_white_rule", "id")),
		}, {
			ResourceName:      "tencentcloud_waf_attack_white_rule.waf_attack_white_rule",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccWafAttackWhiteRule = `

resource "tencentcloud_waf_attack_white_rule" "waf_attack_white_rule" {
  rules = {
  }
}
`
