package tcmg_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudMonitorGrafanaDnsConfigResource_basic -v
func TestAccTencentCloudMonitorGrafanaDnsConfigResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorGrafanaDnsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_dns_config.grafana_dns_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_dns_config.grafana_dns_config", "name_servers.#", "3"),
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
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_dns_config.grafana_dns_config", "name_servers.#", "0"),
				),
			},
		},
	})
}

const testAccMonitorGrafanaDnsConfigVar = `
variable "instance_id" {
  default = "` + tcacctest.DefaultGrafanaInstanceId + `"
}
`

const testAccMonitorGrafanaDnsConfig = testAccMonitorGrafanaDnsConfigVar + `

resource "tencentcloud_monitor_grafana_dns_config" "grafana_dns_config" {
  instance_id  = var.instance_id
  name_servers = ["10.1.2.1", "10.1.2.2", "10.1.2.3"]
}

`

const testAccMonitorGrafanaDnsConfigUp = testAccMonitorGrafanaDnsConfigVar + `

resource "tencentcloud_monitor_grafana_dns_config" "grafana_dns_config" {
  instance_id  = var.instance_id
  name_servers = []
}

`
