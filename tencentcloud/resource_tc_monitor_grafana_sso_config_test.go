package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMonitorGrafanaSsoConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorGrafanaSsoConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_sso_config.grafana_sso_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_monitor_grafana_sso_config.grafana_sso_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMonitorGrafanaSsoConfig = `

resource "tencentcloud_monitor_grafana_sso_config" "grafana_sso_config" {
  instance_id = ""
  enable_s_s_o = 
}

`
