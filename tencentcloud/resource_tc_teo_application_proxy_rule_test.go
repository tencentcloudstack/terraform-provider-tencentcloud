package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixTeoApplicationProxyRule_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoApplicationProxyRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_application_proxy_rule.application_proxy_rule", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_application_proxy_rule.application_proxy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoApplicationProxyRule = `

resource "tencentcloud_teo_application_proxy_rule" "application_proxy_rule" {
  zone_id  = tencentcloud_teo_zone.zone.id
  proxy_id = tencentcloud_teo_application_proxy.application_proxy_rule.proxy_id

  forward_client_ip = "TOA"
  origin_type       = "custom"
  origin_value      = [
    "1.1.1.1:80",
  ]
  port = [
    "80",
  ]
  proto           = "TCP"
  session_persist = false
}

`
