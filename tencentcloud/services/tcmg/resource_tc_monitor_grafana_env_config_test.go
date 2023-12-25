package tcmg_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudMonitorGrafanaEnvConfigResource_basic -v
func TestAccTencentCloudMonitorGrafanaEnvConfigResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorGrafanaEnvConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_env_config.grafana_env_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_env_config.grafana_env_config", "envs.aaa", "ccc"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_env_config.grafana_env_config", "envs.cc", "ccc"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_env_config.grafana_env_config", "envs.w", "ccc"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_grafana_env_config.grafana_env_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMonitorGrafanaEnvConfigUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_env_config.grafana_env_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_env_config.grafana_env_config", "envs.aaa", "ccc"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_env_config.grafana_env_config", "envs.cc", "ccc"),
				),
			},
			{
				Config: testAccMonitorGrafanaEnvConfigNull,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_env_config.grafana_env_config", "id"),
				),
			},
		},
	})
}

const testAccMonitorGrafanaEnvConfigVar = `
variable "instance_id" {
  default = "` + tcacctest.DefaultGrafanaInstanceId + `"
}
`
const testAccMonitorGrafanaEnvConfig = testAccMonitorGrafanaEnvConfigVar + `

resource "tencentcloud_monitor_grafana_env_config" "grafana_env_config" {
  instance_id = var.instance_id
  envs = {
    "aaa" = "ccc"
    "w"   = "ccc"
    "cc"  = "ccc"
  }
}

`

const testAccMonitorGrafanaEnvConfigUp = testAccMonitorGrafanaEnvConfigVar + `

resource "tencentcloud_monitor_grafana_env_config" "grafana_env_config" {
  instance_id = var.instance_id
  envs = {
    "aaa" = "ccc"
    "cc"  = "ccc"
  }
}

`

const testAccMonitorGrafanaEnvConfigNull = testAccMonitorGrafanaEnvConfigVar + `

resource "tencentcloud_monitor_grafana_env_config" "grafana_env_config" {
  instance_id = var.instance_id
}

`
