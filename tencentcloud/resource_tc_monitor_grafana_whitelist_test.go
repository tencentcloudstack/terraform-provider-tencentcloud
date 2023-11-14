package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMonitorGrafanaWhitelistResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorGrafanaWhitelist,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_whitelist.grafana_whitelist", "id")),
			},
			{
				ResourceName:      "tencentcloud_monitor_grafana_whitelist.grafana_whitelist",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMonitorGrafanaWhitelist = `

resource "tencentcloud_monitor_grafana_whitelist" "grafana_whitelist" {
  instance_id = "grafana-abcdefgh"
  whitelist = 
}

`
