package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMysqlSwitchForUpgradeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlSwitchForUpgrade,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mysql_switch_for_upgrade.switch_for_upgrade", "id")),
			},
			{
				ResourceName:      "tencentcloud_mysql_switch_for_upgrade.switch_for_upgrade",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMysqlSwitchForUpgrade = `

resource "tencentcloud_mysql_switch_for_upgrade" "switch_for_upgrade" {
  instance_id = ""
}

`
