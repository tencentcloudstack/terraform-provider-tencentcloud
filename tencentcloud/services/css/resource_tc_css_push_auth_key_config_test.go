package css_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCssPushAuthKeyConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCssPushAuthKeyConfig, tcacctest.DefaultCSSPushDomainName),
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_css_push_auth_key_config.push_auth_key_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_push_auth_key_config.push_auth_key_config", "domain_name", tcacctest.DefaultCSSPushDomainName),
					resource.TestCheckResourceAttr("tencentcloud_css_push_auth_key_config.push_auth_key_config", "enable", "0"),
					resource.TestCheckResourceAttr("tencentcloud_css_push_auth_key_config.push_auth_key_config", "master_auth_key", "test000masterkey"),
					resource.TestCheckResourceAttr("tencentcloud_css_push_auth_key_config.push_auth_key_config", "backup_auth_key", "test000backkey"),
					resource.TestCheckResourceAttr("tencentcloud_css_push_auth_key_config.push_auth_key_config", "auth_delta", "1800"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCssPushAuthKeyConfig_update, tcacctest.DefaultCSSPushDomainName),
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_css_push_auth_key_config.push_auth_key_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_push_auth_key_config.push_auth_key_config", "domain_name", tcacctest.DefaultCSSPushDomainName),
					resource.TestCheckResourceAttr("tencentcloud_css_push_auth_key_config.push_auth_key_config", "enable", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_push_auth_key_config.push_auth_key_config", "master_auth_key", "test000masterkey000updated"),
					resource.TestCheckResourceAttr("tencentcloud_css_push_auth_key_config.push_auth_key_config", "backup_auth_key", "test000backkey000updated"),
					resource.TestCheckResourceAttr("tencentcloud_css_push_auth_key_config.push_auth_key_config", "auth_delta", "3600"),
				),
			},
			{
				ResourceName:      "tencentcloud_css_push_auth_key_config.push_auth_key_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCssPushAuthKeyConfig = `

resource "tencentcloud_css_push_auth_key_config" "push_auth_key_config" {
  domain_name = "%s"
  enable = 0
  master_auth_key = "test000masterkey"
  backup_auth_key = "test000backkey"
  auth_delta = 1800
}

`

const testAccCssPushAuthKeyConfig_update = `

resource "tencentcloud_css_push_auth_key_config" "push_auth_key_config" {
  domain_name = "%s"
  enable = 1
  master_auth_key = "test000masterkey000updated"
  backup_auth_key = "test000backkey000updated"
  auth_delta = 3600
}

`
