package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudMonitorGrafanaWhitelistConfigResource_basic -v
func TestAccTencentCloudMonitorGrafanaWhitelistConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorGrafanaWhitelistConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_whitelist_config.grafana_whitelist_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_whitelist_config.grafana_whitelist_config", "whitelist.#", "3"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_grafana_whitelist_config.grafana_whitelist_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMonitorGrafanaWhitelistConfigUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_whitelist_config.grafana_whitelist_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_whitelist_config.grafana_whitelist_config", "whitelist.#", "2"),
				),
			},
			{
				Config: testAccMonitorGrafanaWhitelistConfigNull,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_whitelist_config.grafana_whitelist_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_whitelist_config.grafana_whitelist_config", "whitelist.#", "0"),
				),
			},
		},
	})
}

const testAccMonitorGrafanaWhitelistConfigVar = `
variable "instance_id" {
  default = "` + defaultGrafanaInstanceId + `"
}
`

const testAccMonitorGrafanaWhitelistConfig = testAccMonitorGrafanaWhitelistConfigVar + `

resource "tencentcloud_monitor_grafana_whitelist_config" "grafana_whitelist_config" {
  instance_id = var.instance_id
  whitelist   = ["10.1.1.1", "10.1.1.2", "10.1.1.3"]
}

`

const testAccMonitorGrafanaWhitelistConfigUp = testAccMonitorGrafanaWhitelistConfigVar + `

resource "tencentcloud_monitor_grafana_whitelist_config" "grafana_whitelist_config" {
  instance_id = var.instance_id
  whitelist   = ["10.1.1.1", "10.1.1.2"]
}

`

const testAccMonitorGrafanaWhitelistConfigNull = testAccMonitorGrafanaWhitelistConfigVar + `

resource "tencentcloud_monitor_grafana_whitelist_config" "grafana_whitelist_config" {
  instance_id = var.instance_id
}

`
