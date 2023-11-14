package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMonitorGrafanaSsoCamConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorGrafanaSsoCamConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_sso_cam_config.grafana_sso_cam_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_monitor_grafana_sso_cam_config.grafana_sso_cam_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMonitorGrafanaSsoCamConfig = `

resource "tencentcloud_monitor_grafana_sso_cam_config" "grafana_sso_cam_config" {
  instance_id = ""
  enable_s_s_o_cam_check = 
}

`
