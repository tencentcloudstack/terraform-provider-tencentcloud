package tcmg_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudMonitorGrafanaSsoCamConfigResource_basic -v
func TestAccTencentCloudMonitorGrafanaSsoCamConfigResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorGrafanaSsoCamConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_sso_cam_config.grafana_sso_cam_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_sso_cam_config.grafana_sso_cam_config", "enable_sso_cam_check", "true"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_grafana_sso_cam_config.grafana_sso_cam_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMonitorGrafanaSsoCamConfigUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_sso_cam_config.grafana_sso_cam_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_sso_cam_config.grafana_sso_cam_config", "enable_sso_cam_check", "false"),
				),
			},
		},
	})
}

const testAccMonitorGrafanaSsoCamConfigVar = `
variable "instance_id" {
  default = "` + tcacctest.DefaultGrafanaInstanceId + `"
}
`

const testAccMonitorGrafanaSsoCamConfig = testAccMonitorGrafanaSsoCamConfigVar + `

resource "tencentcloud_monitor_grafana_sso_cam_config" "grafana_sso_cam_config" {
  instance_id          = var.instance_id
  enable_sso_cam_check = true
}

`

const testAccMonitorGrafanaSsoCamConfigUp = testAccMonitorGrafanaSsoCamConfigVar + `

resource "tencentcloud_monitor_grafana_sso_cam_config" "grafana_sso_cam_config" {
  instance_id          = var.instance_id
  enable_sso_cam_check = false
}

`
