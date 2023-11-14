package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMonitorGrafanaEnvResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorGrafanaEnv,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_env.grafana_env", "id")),
			},
			{
				ResourceName:      "tencentcloud_monitor_grafana_env.grafana_env",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMonitorGrafanaEnv = `

resource "tencentcloud_monitor_grafana_env" "grafana_env" {
  instance_id = "grafana-12345678"
  envs = ""
}

`
