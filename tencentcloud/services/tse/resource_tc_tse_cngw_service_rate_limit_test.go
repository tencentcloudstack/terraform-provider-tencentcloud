package tse_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctse "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tse"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTseCngwServiceRateLimitResource_basic -v
func TestAccTencentCloudTseCngwServiceRateLimitResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTseCngwServiceRateLimitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTseCngwServiceRateLimit,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTseCngwServiceRateLimitExists("tencentcloud_tse_cngw_service_rate_limit.cngw_service_rate_limit"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_service_rate_limit.cngw_service_rate_limit", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service_rate_limit.cngw_service_rate_limit", "gateway_id", tcacctest.DefaultTseGatewayId),
					// resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service_rate_limit.cngw_service_rate_limit", "name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service_rate_limit.cngw_service_rate_limit", "limit_detail.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service_rate_limit.cngw_service_rate_limit", "limit_detail.0.enabled", "true"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service_rate_limit.cngw_service_rate_limit", "limit_detail.0.header", "req"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service_rate_limit.cngw_service_rate_limit", "limit_detail.0.hide_client_headers", "true"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service_rate_limit.cngw_service_rate_limit", "limit_detail.0.is_delay", "true"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service_rate_limit.cngw_service_rate_limit", "limit_detail.0.limit_by", "header"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service_rate_limit.cngw_service_rate_limit", "limit_detail.0.line_up_time", "15"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service_rate_limit.cngw_service_rate_limit", "limit_detail.0.policy", "redis"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service_rate_limit.cngw_service_rate_limit", "limit_detail.0.response_type", "default"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service_rate_limit.cngw_service_rate_limit", "limit_detail.0.qps_thresholds.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service_rate_limit.cngw_service_rate_limit", "limit_detail.0.qps_thresholds.0.max", "100"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service_rate_limit.cngw_service_rate_limit", "limit_detail.0.qps_thresholds.0.unit", "hour"),
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
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctse.NewTseService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tse_cngw_service_rate_limit" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		gatewayId := idSplit[0]
		name := idSplit[1]

		res, err := service.DescribeTseCngwServiceRateLimitById(ctx, gatewayId, name)
		if err != nil {
			return err
		}

		if res != nil && res.Enabled != nil {
			return fmt.Errorf("tse cngwServiceRateLimit %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTseCngwServiceRateLimitExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		gatewayId := idSplit[0]
		name := idSplit[1]

		service := svctse.NewTseService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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

const testAccTseCngwServiceRateLimit = tcacctest.DefaultTseVar + `

resource "tencentcloud_tse_cngw_service_rate_limit" "cngw_service_rate_limit" {
    gateway_id = var.gateway_id
    name       = "b6017eaf-2363-481e-9e93-8d65aaf498cd"

    limit_detail {
        enabled             = true
        header              = "req"
        hide_client_headers = true
        is_delay            = true
        limit_by            = "header"
        line_up_time        = 15
        policy              = "redis"
        response_type       = "default"

        qps_thresholds {
            max  = 100
            unit = "hour"
        }
    }
}

`
