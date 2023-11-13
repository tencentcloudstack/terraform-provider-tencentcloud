package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMonitorSsoAccountResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorSsoAccount,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_monitor_sso_account.sso_account", "id")),
			},
			{
				ResourceName:      "tencentcloud_monitor_sso_account.sso_account",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMonitorSsoAccount = `

resource "tencentcloud_monitor_sso_account" "sso_account" {
  instance_id = &lt;nil&gt;
  user_id = &lt;nil&gt;
  notes = &lt;nil&gt;
  role {
		organization = &lt;nil&gt;
		role = &lt;nil&gt;

  }
}

`
