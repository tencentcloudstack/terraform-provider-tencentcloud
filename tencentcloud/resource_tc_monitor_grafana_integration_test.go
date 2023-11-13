package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMonitorGrafanaIntegrationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorGrafanaIntegration,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_integration.grafana_integration", "id")),
			},
			{
				ResourceName:      "tencentcloud_monitor_grafana_integration.grafana_integration",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMonitorGrafanaIntegration = `

resource "tencentcloud_monitor_grafana_integration" "grafana_integration" {
  instance_id = &lt;nil&gt;
    kind = &lt;nil&gt;
  content = &lt;nil&gt;
}

`
