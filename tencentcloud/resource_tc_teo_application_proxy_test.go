package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTeoApplicationProxy_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoApplicationProxy,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_application_proxy.applicationProxy", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_application_proxy.applicationProxy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoApplicationProxy = `

resource "tencentcloud_teo_application_proxy" "applicationProxy" {
  zone_id              = ""
  zone_name            = ""
  proxy_name           = ""
  plat_type            = ""
  security_type        = ""
  accelerate_type      = ""
  session_persist_time = ""
  proxy_type           = ""
}

`
