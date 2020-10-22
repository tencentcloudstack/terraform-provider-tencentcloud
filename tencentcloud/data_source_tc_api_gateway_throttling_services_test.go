package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAPIGatewaythrottlingServiceDataSourceName = "data.tencentcloud_api_gateway_throttling_services"

func TestAccTencentAPIGatewayThrottlingServicesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckThrottlingServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTestAccTencentAPIGatewayThrottlingServices(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckThrottlingServiceExists("tencentcloud_api_gateway_throttling_service.service"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingServiceDataSourceName+".id", "list.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingServiceDataSourceName+".id", "list.0.service_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingServiceDataSourceName+".id", "list.0.environments.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingServiceDataSourceName+".id", "list.0.environments.0.environment_name"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingServiceDataSourceName+".id", "list.0.environments.0.url"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingServiceDataSourceName+".id", "list.0.environments.0.status"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingServiceDataSourceName+".id", "list.0.environments.0.strategy"),
				),
			},
		},
	})
}

func testAccCheckThrottlingServiceDestroy(s *terraform.State) error {
	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		throttlingService = APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_api_gateway_throttling_service" {
			continue
		}

		serviceId := rs.Primary.Attributes["service_id"]
		environmentNames := rs.Primary.Attributes["environment_names"]
		environmentList, err := throttlingService.DescribeServiceEnvironmentStrategyList(ctx, serviceId)
		if err != nil {
			return err
		}

		for _, v := range environmentList {
			if v == nil || !strings.Contains(environmentNames, *v.EnvironmentName) {
				continue
			}
			if *v.Strategy == STRATEGY || *v.Strategy == STRATEGY_MAX {
				continue
			}

			return fmt.Errorf("throttling service still not restore: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckThrottlingServiceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var (
			logId             = getLogId(contextNil)
			ctx               = context.WithValue(context.TODO(), logIdKey, logId)
			throttlingService = APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("API gateway throttling service %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("API gateway throttling service id is not set")
		}

		serviceId := rs.Primary.Attributes["service_id"]
		environmentNames := rs.Primary.Attributes["environment_names"]
		environmentList, err := throttlingService.DescribeServiceEnvironmentStrategyList(ctx, serviceId)
		if err != nil {
			return err
		}

		for _, v := range environmentList {
			if v == nil || !strings.Contains(environmentNames, *v.EnvironmentName) {
				continue
			}
			if *v.Strategy == STRATEGY {
				return fmt.Errorf("throttling service still not set value: %s", rs.Primary.ID)
			}
		}
		return nil
	}
}

func testAccTestAccTencentAPIGatewayThrottlingServices() string {
	return `
resource "tencentcloud_api_gateway_service" "service" {
  	service_name = "niceservice"
  	protocol     = "http&https"
  	service_desc = "your nice service"
  	net_type     = ["INNER", "OUTER"]
  	ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_throttling_service" "service" {
	service_id        = tencentcloud_api_gateway_service.service.id 
	strategy          = "400"
	environment_names = ["release"]
}

data "tencentcloud_api_gateway_throttling_services" "id" {
    service_id = tencentcloud_api_gateway_throttling_service.service.service_id
}
`
}
