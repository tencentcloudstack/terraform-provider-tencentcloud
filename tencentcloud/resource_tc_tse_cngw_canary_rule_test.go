package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudTseCngwCanaryRuleResource_basic -v
func TestAccTencentCloudTseCngwCanaryRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTseCngwCanaryRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTseCngwCanaryRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTseCngwCanaryRuleExists("tencentcloud_tse_cngw_canary_rule.cngw_canary_rule"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_canary_rule.cngw_canary_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_canary_rule.cngw_canary_rule", "gateway_id", defaultTseGatewayId),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_canary_rule.cngw_canary_rule", "tags.created", "terraform"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_canary_rule.cngw_canary_rule", "canary_rule.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_canary_rule.cngw_canary_rule", "canary_rule.0.enabled", "true"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_canary_rule.cngw_canary_rule", "canary_rule.0.priority", "100"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_canary_rule.cngw_canary_rule", "canary_rule.0.balanced_service_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_canary_rule.cngw_canary_rule", "canary_rule.0.balanced_service_list.0.percent", "100"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_canary_rule.cngw_canary_rule", "canary_rule.0.balanced_service_list.0.service_name", "terraform-test-canary_rule"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_canary_rule.cngw_canary_rule", "canary_rule.0.balanced_service_list.0.service_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_canary_rule.cngw_canary_rule", "canary_rule.0.balanced_service_list.0.upstream_name"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_canary_rule.cngw_canary_rule", "canary_rule.0.condition_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_canary_rule.cngw_canary_rule", "canary_rule.0.condition_list.0.key", "test"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_canary_rule.cngw_canary_rule", "canary_rule.0.condition_list.0.operator", "eq"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_canary_rule.cngw_canary_rule", "canary_rule.0.condition_list.0.type", "query"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_canary_rule.cngw_canary_rule", "canary_rule.0.condition_list.0.value", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_tse_cngw_canary_rule.cngw_canary_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTseCngwCanaryRuleDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TseService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tse_cngw_canary_rule" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		gatewayId := idSplit[0]
		serviceId := idSplit[1]
		priority := idSplit[2]

		res, err := service.DescribeTseCngwCanaryRuleById(ctx, gatewayId, serviceId, priority)
		if err != nil {
			if sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkErr.Code == "FailedOperation.FailedOperation" {
					return nil
				}
			}
			return err
		}

		if res != nil {
			return fmt.Errorf("tse cngwCanaryRule %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTseCngwCanaryRuleExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		gatewayId := idSplit[0]
		serviceId := idSplit[1]
		priority := idSplit[2]

		service := TseService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTseCngwCanaryRuleById(ctx, gatewayId, serviceId, priority)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tse cngwCanaryRule %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTseCngwCanaryRule = `

resource "tencentcloud_tse_cngw_service" "cngw_service" {
	gateway_id = "gateway-ddbb709b"
	name       = "terraform-test-canary_rule"
	path       = "/test"
	protocol   = "http"
	retries    = 5
	tags = {
	  "created" = "terraform"
	}
	timeout       = 6000
	upstream_type = "IPList"
  
	upstream_info {
	  algorithm                   = "round-robin"
	  auto_scaling_cvm_port       = 80
	  auto_scaling_group_id       = "asg-519acdug"
	  auto_scaling_hook_status    = "Normal"
	  auto_scaling_tat_cmd_status = "Normal"
	  port                        = 0
	  slow_start                  = 20
  
	  targets {
		host   = "192.168.0.1"
		port   = 80
		weight = 100
	  }
	}
}
  
resource "tencentcloud_tse_cngw_canary_rule" "cngw_canary_rule" {
	gateway_id = tencentcloud_tse_cngw_service.cngw_service.gateway_id
	service_id = tencentcloud_tse_cngw_service.cngw_service.service_id
	tags       = {
	  "created" = "terraform"
	}
  
	canary_rule {
	  enabled  = true
	  priority = 100
  
	  balanced_service_list {
		percent       = 100
		service_id    = tencentcloud_tse_cngw_service.cngw_service.service_id
		service_name  = tencentcloud_tse_cngw_service.cngw_service.name
	  }
  
	  condition_list {
		key      = "test"
		operator = "eq"
		type     = "query"
		value    = "1"
	  }
	}
}

`
