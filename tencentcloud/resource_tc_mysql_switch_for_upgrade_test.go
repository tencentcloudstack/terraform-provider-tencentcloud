package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixMysqlSwitchForUpgradeResource_basic(t *testing.T) {
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
		},
	})
}

const testAccMysqlSwitchForUpgrade = `

resource "tencentcloud_mysql_switch_for_upgrade" "switch_for_upgrade" {
	instance_id = "cdb-d9gbh7lt"
}

`
