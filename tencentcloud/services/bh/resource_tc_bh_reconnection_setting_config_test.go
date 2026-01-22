package bh_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudBhReconnectionSettingConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBhReconnectionSettingConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bh_reconnection_setting_config.example", "id"),
				),
			},
			{
				Config: testAccBhReconnectionSettingConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bh_reconnection_setting_config.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_bh_reconnection_setting_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccBhReconnectionSettingConfig = `
resource "tencentcloud_bh_reconnection_setting_config" "example" {
  reconnection_max_count = 5
  enable                 = false
}
`

const testAccBhReconnectionSettingConfigUpdate = `
resource "tencentcloud_bh_reconnection_setting_config" "example" {
  reconnection_max_count = 3
  enable                 = true
}
`
