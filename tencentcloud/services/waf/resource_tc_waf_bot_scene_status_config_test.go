package waf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWafBotSceneUCBRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafBotSceneUCBRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_scene_status_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_scene_status_config.example", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_scene_status_config.example", "scene_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_scene_status_config.example", "status"),
				),
			},
			{
				Config: testAccWafBotSceneUCBRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_scene_status_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_scene_status_config.example", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_scene_status_config.example", "scene_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_scene_status_config.example", "status"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_bot_scene_status_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWafBotSceneUCBRule = `
resource "tencentcloud_waf_bot_scene_status_config" "example" {
  domain   = "example.com"
  scene_id = "3024324123"
  status   = true
}
`

const testAccWafBotSceneUCBRuleUpdate = `
resource "tencentcloud_waf_bot_scene_status_config" "example" {
  domain   = "example.com"
  scene_id = "3024324123"
  status   = false
}
`
