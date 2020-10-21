package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudAPIGateWayThrottlingService(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckThrottlingServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccThrottlingService,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckThrottlingServiceExists("tencentcloud_api_gateway_throttling_service.foo"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_throttling_service.foo", "service_id"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_throttling_service.foo", "strategy", "400"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_throttling_service.foo", "environment_names.#", "1"),
				),
			},
			{
				Config: testAccThrottlingServiceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckThrottlingServiceExists("tencentcloud_api_gateway_throttling_service.foo"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_throttling_service.foo", "service_id"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_throttling_service.foo", "strategy", "400"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_throttling_service.foo", "environment_names.#", "2"),
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

const testServiceBase = `
resource "tencentcloud_api_gateway_service" "service" {
  	service_name = "niceservice"
  	protocol     = "http&https"
  	service_desc = "your nice service"
  	net_type     = ["INNER", "OUTER"]
  	ip_version   = "IPv4"
}
`

const testAccThrottlingService = testServiceBase + `
resource "tencentcloud_api_gateway_throttling_service" "foo" {
	service_id        = tencentcloud_api_gateway_service.service.id 
	strategy          = "400"
	environment_names = ["release"]
}
`

const testAccThrottlingServiceUpdate = testServiceBase + `
resource "tencentcloud_api_gateway_throttling_service" "foo" {
	service_id        = tencentcloud_api_gateway_service.service.id
	strategy          = "400"
	environment_names = ["release", "test"]
}
`
