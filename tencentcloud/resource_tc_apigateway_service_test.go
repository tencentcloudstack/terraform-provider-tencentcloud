package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudApigatewayServiceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApigatewayService,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_apigateway_service.service", "id")),
			},
			{
				ResourceName:      "tencentcloud_apigateway_service.service",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccApigatewayService = `

resource "tencentcloud_apigateway_service" "service" {
  service_name = ""
  protocol = ""
  service_desc = ""
  net_types = 
  ip_version = ""
  set_server_name = ""
  app_id_type = ""
  instance_id = ""
  uniq_vpc_id = ""
  tags = {
    "createdBy" = "terraform"
  }
}

`
