package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbSwitchForUpgradeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbSwitchForUpgrade,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdb_switch_for_upgrade.switch_for_upgrade", "id")),
			},
			{
				ResourceName:      "tencentcloud_cdb_switch_for_upgrade.switch_for_upgrade",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdbSwitchForUpgrade = `

resource "tencentcloud_cdb_switch_for_upgrade" "switch_for_upgrade" {
  instance_id = ""
}

`
