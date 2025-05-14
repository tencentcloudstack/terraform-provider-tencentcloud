package waf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWafBotSceneStatusConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafBotSceneStatusConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_scene_status_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_scene_status_config.example", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_scene_status_config.example", "scene_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_scene_status_config.example", "status"),
				),
			},
			{
				Config: testAccWafBotSceneStatusConfigUpdate,
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

const testAccWafBotSceneStatusConfig = `
resource "tencentcloud_waf_bot_scene_status_config" "example" {
  domain   = "example.com"
  scene_id = "3024324123"
  status   = true
}
`

const testAccWafBotSceneStatusConfigUpdate = `
resource "tencentcloud_waf_bot_scene_status_config" "example" {
  domain   = "example.com"
  scene_id = "3024324123"
  status   = false
}
`
