package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudMonitorGrafanaSsoConfigResource_basic -v
func TestAccTencentCloudMonitorGrafanaSsoConfigResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorGrafanaSsoConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_sso_config.grafana_sso_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_sso_config.grafana_sso_config", "enable_sso", "true"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_grafana_sso_config.grafana_sso_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMonitorGrafanaSsoConfigUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_sso_config.grafana_sso_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_sso_config.grafana_sso_config", "enable_sso", "false"),
				),
			},
		},
	})
}

const testAccMonitorGrafanaSsoConfigVar = `
variable "instance_id" {
  default = "` + defaultGrafanaInstanceId + `"
}
`

const testAccMonitorGrafanaSsoConfig = testAccMonitorGrafanaSsoConfigVar + `

resource "tencentcloud_monitor_grafana_sso_config" "grafana_sso_config" {
  instance_id = var.instance_id
  enable_sso  = true
}

`

const testAccMonitorGrafanaSsoConfigUp = testAccMonitorGrafanaSsoConfigVar + `

resource "tencentcloud_monitor_grafana_sso_config" "grafana_sso_config" {
  instance_id = var.instance_id
  enable_sso  = false
}

`
