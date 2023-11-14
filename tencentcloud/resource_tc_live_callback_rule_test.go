package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLiveCallbackRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveCallbackRule,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_live_callback_rule.callback_rule", "id")),
			},
			{
				ResourceName:      "tencentcloud_live_callback_rule.callback_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLiveCallbackRule = `

resource "tencentcloud_live_callback_rule" "callback_rule" {
  domain_name = "5000.livepush.myqcloud.com"
  app_name = "live"
  template_id = 1000
}

`
