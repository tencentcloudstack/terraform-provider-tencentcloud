package cfw_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCfwIpsModeSwitchResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwIpsModeSwitch,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_ips_mode_switch.example", "id"),
				),
			},
			{
				Config: testAccCfwIpsModeSwitchUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_ips_mode_switch.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_cfw_ips_mode_switch.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCfwIpsModeSwitch = `
resource "tencentcloud_cfw_ips_mode_switch" "example" {
  mode = 1
}
`

const testAccCfwIpsModeSwitchUpdate = `
resource "tencentcloud_cfw_ips_mode_switch" "example" {
  mode = 0
}
`
