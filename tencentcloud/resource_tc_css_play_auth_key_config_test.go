package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudCssPlayAuthKeyConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssPlayAuthKeyConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_css_play_auth_key_config.play_auth_key_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_css_play_auth_key_config.play_auth_key_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCssPlayAuthKeyConfig = `

resource "tencentcloud_css_play_auth_key_config" "play_auth_key_config" {
  domain_name = "5000.livepush.myqcloud.com"
  enable = 1
  auth_key = "xx"
  auth_delta = 60
  auth_back_key = "xx"
}

`
