package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTseCngwServiceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseCngwService,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_service.cngw_service", "id")),
			},
			{
				ResourceName:      "tencentcloud_tse_cngw_service.cngw_service",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTseCngwService = `

resource "tencentcloud_tse_cngw_service" "cngw_service" {
  gateway_id = "gateway-xxxxxx"
  name = "serviceA"
  protocol = "https"
  path = "/test"
  timeout = 3000
  retries = 10
  upstream_type = "IPList"
  upstream_info {
		host = "123.123.123.123"
		port = 33
		source_i_d = "ins-xxxxxx"
		namespace = "test"
		service_name = "orderService"
		targets {
			host = "123.123.123.123"
			port = 80
			weight = 10
			health = "healthy"
			created_time = ""
			source = ""
		}
		source_type = ""
		scf_type = ""
		scf_namespace = ""
		scf_lambda_name = ""
		scf_lambda_qualifier = ""
		slow_start = 
		algorithm = ""
		auto_scaling_group_i_d = ""
		auto_scaling_cvm_port = 
		auto_scaling_tat_cmd_status = ""
		auto_scaling_hook_status = ""
		source_name = ""
		real_source_type = ""

  }
}

`
