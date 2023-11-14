package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMonitorGrafanaPluginResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorGrafanaPlugin,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_plugin.grafana_plugin", "id")),
			},
			{
				ResourceName:      "tencentcloud_monitor_grafana_plugin.grafana_plugin",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMonitorGrafanaPlugin = `

resource "tencentcloud_monitor_grafana_plugin" "grafana_plugin" {
  instance_id = &lt;nil&gt;
  plugin_id = &lt;nil&gt;
  version = &lt;nil&gt;
}

`
