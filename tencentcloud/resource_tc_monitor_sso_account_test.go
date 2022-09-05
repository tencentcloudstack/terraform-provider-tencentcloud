package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudMonitorSsoAccount_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorSsoAccount,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_sso_account.ssoAccount", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_sso_account.ssoAccount",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMonitorSsoAccount = `

resource "tencentcloud_monitor_sso_account" "ssoAccount" {
  instance_id = ""
  user_id = ""
  notes = ""
  role {
			organization = ""
			role = ""

  }
}

`
