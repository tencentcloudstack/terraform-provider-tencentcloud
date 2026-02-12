package bh_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudBhAuthModeSettingConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccBhAuthModeSettingConfig,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_bh_auth_mode_setting_config.bh_auth_mode_setting_config", "id")),
		}, {
			ResourceName:      "tencentcloud_bh_auth_mode_setting_config.bh_auth_mode_setting_config",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccBhAuthModeSettingConfig = `

resource "tencentcloud_bh_auth_mode_setting_config" "bh_auth_mode_setting_config" {
}
`
