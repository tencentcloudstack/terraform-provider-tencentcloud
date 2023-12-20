package css_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCssPlayAuthKeyConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCssPlayAuthKeyConfig, tcacctest.DefaultCSSPlayDomainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_play_auth_key_config.play_auth_key_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_play_auth_key_config.play_auth_key_config", "domain_name", tcacctest.DefaultCSSPlayDomainName),
					resource.TestCheckResourceAttr("tencentcloud_css_play_auth_key_config.play_auth_key_config", "enable", "0"),
					resource.TestCheckResourceAttr("tencentcloud_css_play_auth_key_config.play_auth_key_config", "auth_key", "test000key"),
					resource.TestCheckResourceAttr("tencentcloud_css_play_auth_key_config.play_auth_key_config", "auth_back_key", "test000backkey"),
					resource.TestCheckResourceAttr("tencentcloud_css_play_auth_key_config.play_auth_key_config", "auth_delta", "1800"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCssPlayAuthKeyConfig_update, tcacctest.DefaultCSSPlayDomainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_play_auth_key_config.play_auth_key_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_play_auth_key_config.play_auth_key_config", "domain_name", tcacctest.DefaultCSSPlayDomainName),
					resource.TestCheckResourceAttr("tencentcloud_css_play_auth_key_config.play_auth_key_config", "enable", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_play_auth_key_config.play_auth_key_config", "auth_key", "test000key000updated"),
					resource.TestCheckResourceAttr("tencentcloud_css_play_auth_key_config.play_auth_key_config", "auth_back_key", "test000backkey000updated"),
					resource.TestCheckResourceAttr("tencentcloud_css_play_auth_key_config.play_auth_key_config", "auth_delta", "3600"),
				),
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
  domain_name = "%s"
  enable = 0
  auth_key = "test000key" 
  auth_delta = 1800
  auth_back_key = "test000backkey"
}

`

const testAccCssPlayAuthKeyConfig_update = `

resource "tencentcloud_css_play_auth_key_config" "play_auth_key_config" {
  domain_name = "%s"
  enable = 1
  auth_key = "test000key000updated" 
  auth_delta = 3600
  auth_back_key = "test000backkey000updated"
}

`
