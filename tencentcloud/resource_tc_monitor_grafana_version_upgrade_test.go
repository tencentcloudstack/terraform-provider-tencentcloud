package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMonitorGrafanaVersionUpgradeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorGrafanaVersionUpgrade,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_version_upgrade.grafana_version_upgrade", "id")),
			},
			{
				ResourceName:      "tencentcloud_monitor_grafana_version_upgrade.grafana_version_upgrade",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMonitorGrafanaVersionUpgrade = `

resource "tencentcloud_monitor_grafana_version_upgrade" "grafana_version_upgrade" {
  instance_id = ""
  alias = ""
}

`
