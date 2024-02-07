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
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -test.run TestAccTencentCloudTseCngwStrategyResource_basic -v -timeout=0
func TestAccTencentCloudTseCngwStrategyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTseCngwStrategyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTseCngwStrategy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTseCngwStrategyExists("tencentcloud_tse_cngw_strategy.cngw_strategy"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_strategy.cngw_strategy", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "description", "test"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "strategy_name", "test-cron"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.max_replicas", "2"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_down.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_down.0.select_policy", "Max"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_down.0.stabilization_window_seconds", "300"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_down.0.policies.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_down.0.policies.0.period_seconds", "10"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_down.0.policies.0.type", "Pods"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_down.0.policies.0.value", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_down.0.policies.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_up.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_up.0.select_policy", "Max"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_up.0.stabilization_window_seconds", "30"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_up.0.policies.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_up.0.policies.0.period_seconds", "10"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_up.0.policies.0.type", "Pods"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_up.0.policies.0.value", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.metrics.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.metrics.0.resource_name", "cpu"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.metrics.0.target_value", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.metrics.0.type", "Resource"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "cron_config.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "cron_config.0.params.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "cron_config.0.params.0.crontab", "0 00 00 * * *"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "cron_config.0.params.0.period", "* * *"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "cron_config.0.params.0.start_at", "00:00"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "cron_config.0.params.0.target_replicas", "2"),
				),
			},
			{
				ResourceName:      "tencentcloud_tse_cngw_strategy.cngw_strategy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTseCngwStrategyUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTseCngwStrategyExists("tencentcloud_tse_cngw_strategy.cngw_strategy"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_strategy.cngw_strategy", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "description", "testup"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "strategy_name", "test-cron"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.max_replicas", "2"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_down.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_down.0.select_policy", "Max"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_down.0.stabilization_window_seconds", "301"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_down.0.policies.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_down.0.policies.0.period_seconds", "9"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_down.0.policies.0.type", "Pods"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_down.0.policies.0.value", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_down.0.policies.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_up.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_up.0.select_policy", "Max"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_up.0.stabilization_window_seconds", "31"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_up.0.policies.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_up.0.policies.0.period_seconds", "10"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_up.0.policies.0.type", "Pods"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.behavior.0.scale_up.0.policies.0.value", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.metrics.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.metrics.0.resource_name", "cpu"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.metrics.0.target_value", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "config.0.metrics.0.type", "Resource"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "cron_config.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "cron_config.0.params.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "cron_config.0.params.0.crontab", "0 00 00 * * *"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "cron_config.0.params.0.period", "* * *"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "cron_config.0.params.0.start_at", "00:00"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy.cngw_strategy", "cron_config.0.params.0.target_replicas", "2"),
				),
			},
		},
	})
}

func testAccCheckTseCngwStrategyDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctse.NewTseService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tse_cngw_service" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		gatewayId := idSplit[0]
		strategyId := idSplit[1]

		res, err := service.DescribeTseCngwStrategyById(ctx, gatewayId, strategyId)
		if err != nil {
			if sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkErr.Code == "ResourceNotFound.ResourceNotFound" {
					return nil
				}
			}
			return err
		}

		if res != nil {
			return fmt.Errorf("tse cngwService %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTseCngwStrategyExists(r string) resource.TestCheckFunc {
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
		strategyId := idSplit[1]

		service := svctse.NewTseService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeTseCngwStrategyById(ctx, gatewayId, strategyId)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tse cngwService %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTseCngwStrategy = testAccTseCngwGateway + `

resource "tencentcloud_tse_cngw_strategy" "cngw_strategy" {
  description   = "test"
  gateway_id    = tencentcloud_tse_cngw_gateway.cngw_gateway.id
  strategy_name = "test-cron"

  config {
    max_replicas = 2
    behavior {
      scale_down {
        select_policy                = "Max"
        stabilization_window_seconds = 300

        policies {
          period_seconds = 10
          type           = "Pods"
          value          = 1
        }
      }
      scale_up {
        select_policy                = "Max"
        stabilization_window_seconds = 30

        policies {
          period_seconds = 10
          type           = "Pods"
          value          = 1
        }
      }
    }

    metrics {
      resource_name = "cpu"
      target_value  = 1
      type          = "Resource"
    }
  }

  cron_config {
    params {
      crontab         = "0 00 00 * * *"
      period          = "* * *"
      start_at        = "00:00"
      target_replicas = 2
    }
  }
}

`
const testAccTseCngwStrategyUp = testAccTseCngwGateway + `

resource "tencentcloud_tse_cngw_strategy" "cngw_strategy" {
  description   = "testup"
  gateway_id    = tencentcloud_tse_cngw_gateway.cngw_gateway.id
  strategy_name = "test-cron"

  config {
    max_replicas = 2
    behavior {
      scale_down {
        select_policy                = "Max"
        stabilization_window_seconds = 301

        policies {
          period_seconds = 9
          type           = "Pods"
          value          = 1
        }
      }
      scale_up {
        select_policy                = "Max"
        stabilization_window_seconds = 31

        policies {
          period_seconds = 10
          type           = "Pods"
          value          = 1
        }
      }
    }

    metrics {
      resource_name = "cpu"
      target_value  = 1
      type          = "Resource"
    }
  }

  cron_config {
    params {
      crontab         = "0 00 00 * * *"
      period          = "* * *"
      start_at        = "00:00"
      target_replicas = 2
    }
  }
}

`
