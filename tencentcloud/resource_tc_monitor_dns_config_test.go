package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMonitorDnsConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorDnsConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_monitor_dns_config.dns_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_monitor_dns_config.dns_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMonitorDnsConfig = `

resource "tencentcloud_monitor_dns_config" "dns_config" {
  instance_id = "grafana-12345678"
  name_servers = 
}

`
