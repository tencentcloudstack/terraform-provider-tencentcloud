package teo_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTeoL4proxyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoL4proxy,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_l4proxy.l4proxy", "id")),
			},
			{
				ResourceName:      "tencentcloud_teo_l4proxy.l4proxy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoL4proxy = `

resource "tencentcloud_teo_l4proxy" "l4proxy" {
  zone_id = "zone-21xfqlh4qjee"
  proxy_name = "test-proxy"
  area = "mainland"
  ipv6 = "on"
  static_ip = "on"
  accelerate_mainland = "off"
  d_dos_protection_config {
		level_mainland = "BASE30_MAX300"
		max_bandwidth_mainland = 
		level_overseas = ""

  }
}

`
