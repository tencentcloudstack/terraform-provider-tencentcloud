package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMonitorGrafanaDnsConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorGrafanaDnsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_dns_config.grafana_dns_config", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_grafana_dns_config.grafana_dns_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMonitorGrafanaDnsConfigUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_dns_config.grafana_dns_config", "id"),
				),
			},
		},
	})
}

const testAccMonitorGrafanaDnsConfig = `

resource "tencentcloud_monitor_grafana_dns_config" "grafana_dns_config" {
  instance_id  = "grafana-dp2hnnfa"
  name_servers = ["10.1.2.1", "10.1.2.2", "10.1.2.3"]
}

`

const testAccMonitorGrafanaDnsConfigUp = `

resource "tencentcloud_monitor_grafana_dns_config" "grafana_dns_config" {
  instance_id  = "grafana-dp2hnnfa"
  name_servers = []
}

`
