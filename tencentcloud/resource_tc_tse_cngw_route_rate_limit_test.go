package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTseCngwRouteRateLimitResource_basic -v
func TestAccTencentCloudNeedFixTseCngwRouteRateLimitResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTseCngwRouteRateLimitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTseCngwRouteRateLimit,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTseCngwRouteRateLimitExists("tencentcloud_tse_cngw_route_rate_limit.cngw_route_rate_limit"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_route_rate_limit.cngw_route_rate_limit", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_route_rate_limit.cngw_route_rate_limit", "gateway_id", defaultTseGatewayId),
				),
			},
			{
				ResourceName:      "tencentcloud_tse_cngw_route_rate_limit.cngw_route_rate_limit",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTseCngwRouteRateLimitDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TseService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tse_cngw_route_rate_limit" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		gatewayId := idSplit[0]
		routeId := idSplit[1]

		res, err := service.DescribeTseCngwRouteRateLimitById(ctx, gatewayId, routeId)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tse cngwRouteRateLimit %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTseCngwRouteRateLimitExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		gatewayId := idSplit[0]
		routeId := idSplit[1]

		service := TseService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTseCngwRouteRateLimitById(ctx, gatewayId, routeId)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tse cngwRouteRateLimit %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTseCngwRouteRateLimit = `

resource "tencentcloud_tse_cngw_route_rate_limit" "cngw_route_rate_limit" {
  gateway_id = "gateway-xxxxxx"
  route_id = "gateway-xxxxxx"
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
  tags = {
    "createdBy" = "terraform"
  }
}

`
