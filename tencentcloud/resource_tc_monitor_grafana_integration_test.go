package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudMonitorGrafanaIntegration_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorGrafanaIntegration,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_integration.grafanaIntegration", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_grafana_integration.grafanaIntegration",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMonitorGrafanaIntegration = `

resource "tencentcloud_monitor_grafana_integration" "grafanaIntegration" {
  instance_id = ""
    kind = ""
  content = ""
}

`
