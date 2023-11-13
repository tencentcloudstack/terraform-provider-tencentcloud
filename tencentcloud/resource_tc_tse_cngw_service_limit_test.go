package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTseCngwServiceLimitResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseCngwServiceLimit,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_service_limit.cngw_service_limit", "id")),
			},
			{
				ResourceName:      "tencentcloud_tse_cngw_service_limit.cngw_service_limit",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTseCngwServiceLimit = `

resource "tencentcloud_tse_cngw_service_limit" "cngw_service_limit" {
  gateway_id = "gateway-xxxxxx"
  name = "451a9920-e67a-4519-af41-fccac0e72005"
  limit_detail {
		enabled = true
		qps_thresholds {
			unit = "second"
			max = 50
		}
		limit_by = "ip"
		response_type = "default"
		hide_client_headers = false
		is_delay = false
		path = "/test"
		header = "auth"
		external_redis {
			redis_host = ""
			redis_password = ""
			redis_port = 
			redis_timeout = 
		}
		policy = "redis"
		rate_limit_response {
			body = ""
			headers {
				key = ""
				value = ""
			}
			http_status = 
		}
		rate_limit_response_url = ""
		line_up_time = 

  }
}

`
