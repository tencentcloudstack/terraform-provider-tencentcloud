package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixTeoApplicationProxy_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoApplicationProxy,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_application_proxy.application_proxy", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_application_proxy.application_proxy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoApplicationProxy = `

resource "tencentcloud_teo_application_proxy" "application_proxy" {
  zone_id   = tencentcloud_teo_zone.zone.id
  zone_name = "sfurnace.work"

  accelerate_type      = 1
  security_type        = 1
  plat_type            = "domain"
  proxy_name           = "www.sfurnace.work"
  proxy_type           = "hostname"
  session_persist_time = 2400
}

`
