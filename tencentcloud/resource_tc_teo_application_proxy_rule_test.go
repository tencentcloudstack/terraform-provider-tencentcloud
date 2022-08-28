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
					resource.TestCheckResourceAttrSet("tencentcloud_teo_application_proxy_rule.applicationProxyRule", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_application_proxy_rule.applicationProxyRule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoApplicationProxyRule = `

resource "tencentcloud_teo_application_proxy_rule" "applicationProxyRule" {
  zone_id           = ""
  proxy_id          = ""
  proto             = ""
  port              = ""
  origin_type       = ""
  origin_value      = ""
  forward_client_ip = ""
  session_persist   = ""
}

`
