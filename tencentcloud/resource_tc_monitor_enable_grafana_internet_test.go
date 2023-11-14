package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMonitorEnableGrafanaInternetResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorEnableGrafanaInternet,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_monitor_enable_grafana_internet.enable_grafana_internet", "id")),
			},
			{
				ResourceName:      "tencentcloud_monitor_enable_grafana_internet.enable_grafana_internet",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMonitorEnableGrafanaInternet = `

resource "tencentcloud_monitor_enable_grafana_internet" "enable_grafana_internet" {
  instance_i_d = "grafana-kleu3gt0"
  enable_internet = true
}

`
