package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTseCngwServiceRateLimitResource_basic -v
func TestAccTencentCloudNeedFixTseCngwServiceRateLimitResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTseCngwServiceRateLimitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTseCngwServiceRateLimit,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTseCngwServiceRateLimitExists("tencentcloud_tse_cngw_service_rate_limit.cngw_service_rate_limit"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_service_rate_limit.cngw_service_rate_limit", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service_rate_limit.cngw_service_rate_limit", "gateway_id", defaultTseGatewayId),
				),
			},
			{
				ResourceName:      "tencentcloud_tse_cngw_service_rate_limit.cngw_service_rate_limit",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTseCngwServiceRateLimitDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TseService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tse_cngw_service_rate_limit" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		gatewayId := idSplit[0]
		name := idSplit[1]

		res, err := service.DescribeTseCngwServiceRateLimitById(ctx, gatewayId, name)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tse cngwServiceRateLimit %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTseCngwServiceRateLimitExists(r string) resource.TestCheckFunc {
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
		name := idSplit[1]

		service := TseService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTseCngwServiceRateLimitById(ctx, gatewayId, name)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tse cngwServiceRateLimit %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTseCngwServiceRateLimit = DefaultTseVar + `

resource "tencentcloud_tse_cngw_service_rate_limit" "cngw_service_rate_limit" {
    gateway_id = "gateway-ddbb709b"
    name       = "terraform-test"

    limit_detail {
        enabled             = true
        hide_client_headers = true
        is_delay            = true
        limit_by            = "ip"
        line_up_time        = 11
        policy              = "redis"
        response_type       = "default"

        qps_thresholds {
            max  = 14
            unit = "second"
        }
    }
}

`
