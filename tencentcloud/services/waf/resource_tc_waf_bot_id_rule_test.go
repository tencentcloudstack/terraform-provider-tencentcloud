package waf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWafBotIdRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafBotIdRule,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_id_rule.example", "id")),
			},
			{
				Config: testAccWafBotIdRuleUpdate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_id_rule.example", "id")),
			},
			{
				ResourceName:      "tencentcloud_waf_bot_id_rule.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWafBotIdRule = `
resource "tencentcloud_waf_bot_id_rule" "example" {
  domain        = "demo.com"
  scene_id      = "3000000001"
  protect_level = "normal"
  global_switch = 5
}
`

const testAccWafBotIdRuleUpdate = `
resource "tencentcloud_waf_bot_id_rule" "example" {
  domain        = "demo.com"
  scene_id      = "3000000001"
  protect_level = "normal"
  global_switch = 0
}
`
