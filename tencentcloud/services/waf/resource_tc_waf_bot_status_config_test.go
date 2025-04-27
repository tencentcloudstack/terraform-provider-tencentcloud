package waf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWafBotStatusConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafBotStatusConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_status_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_status_config.example", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_status_config.example", "status"),
				),
			},
			{
				Config: testAccWafBotStatusConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_status_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_status_config.example", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_bot_status_config.example", "status"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_bot_status_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWafBotStatusConfig = `
resource "tencentcloud_waf_bot_status_config" "example" {
  domain      = "example.com"
  status      = "1"
}
`

const testAccWafBotStatusConfigUpdate = `
resource "tencentcloud_waf_bot_status_config" "example" {
  domain      = "example.com"
  status      = "0"
}
`
