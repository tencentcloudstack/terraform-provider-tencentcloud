package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudMonitorGrafanaPlugin_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorGrafanaPlugin,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_plugin.grafanaPlugin", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_grafana_plugin.grafanaPlugin",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMonitorGrafanaPlugin = `

resource "tencentcloud_monitor_grafana_plugin" "grafanaPlugin" {
  instance_id = ""
  plugin_id = ""
  version = ""
}

`
