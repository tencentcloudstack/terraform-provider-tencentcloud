package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudMonitorGrafanaVersionUpgradeResource_basic -v
func TestAccTencentCloudMonitorGrafanaVersionUpgradeResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorGrafanaVersionUpgrade,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_version_upgrade.grafana_version_upgrade", "id"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_version_upgrade.grafana_version_upgrade", "alias", "v8.2.7"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_grafana_version_upgrade.grafana_version_upgrade",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMonitorGrafanaVersionUpgradeVar = `
variable "instance_id" {
  default = "` + defaultGrafanaInstanceId + `"
}
`

const testAccMonitorGrafanaVersionUpgrade = testAccMonitorGrafanaVersionUpgradeVar + `

resource "tencentcloud_monitor_grafana_version_upgrade" "grafana_version_upgrade" {
  instance_id = var.instance_id
  alias       = "v8.2.7"
}

`
