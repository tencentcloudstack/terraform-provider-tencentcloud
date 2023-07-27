package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixTseCngwRouteResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseCngwRoute,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_route.cngw_route", "id")),
			},
			{
				ResourceName:      "tencentcloud_tse_cngw_route.cngw_route",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTseCngwRoute = `

resource "tencentcloud_tse_cngw_route" "cngw_route" {
  gateway_id = "gateway-xxxxxx"
  service_i_d = "451a9920-e67a-4519-af41-fccac0e72005"
  route_name = "routeA"
  methods = 
  hosts = 
  paths = 
  protocols = 
  preserve_host = true
  https_redirect_status_code = 302
  strip_path = true
  force_https = 
  destination_ports = 
  headers {
		key = "token"
		value = "xxxxxx"

  }
  tags = {
    "createdBy" = "terraform"
  }
}

`
