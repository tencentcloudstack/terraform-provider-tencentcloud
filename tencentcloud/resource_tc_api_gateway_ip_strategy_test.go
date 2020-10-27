package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudAPIGateWayIPStrategy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testApiIPStrategyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testApiIPStrategy_basic,
				Check: resource.ComposeTestCheckFunc(
					testApiIPStrategyExists("tencentcloud_api_gateway_ip_strategy.test"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_ip_strategy.test", "service_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_ip_strategy.test", "strategy_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_ip_strategy.test", "strategy_type"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_ip_strategy.test", "strategy_data", "9.9.9.9"),
				),
			},
			{
				ResourceName:      "tencentcloud_api_gateway_ip_strategy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testApiIPStrategy_update,
				Check: resource.ComposeTestCheckFunc(
					testApiIPStrategyExists("tencentcloud_api_gateway_ip_strategy.test"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_ip_strategy.test", "service_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_ip_strategy.test", "strategy_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_ip_strategy.test", "strategy_type"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_ip_strategy.test", "strategy_data", "7.7.7.7"),
				),
			},
		},
	})
}

func testApiIPStrategyDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_api_gateway_ip_strategy" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("IP strategy id is borken, id is %s", rs.Primary.ID)
		}
		serviceId := idSplit[0]
		strategyId := idSplit[1]

		has, err := service.DescribeIPStrategyHas(ctx, serviceId, strategyId)
		if err != nil {
			return err
		}

		if has {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][IP strategy][Destroy] check: IP strategy still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testApiIPStrategyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][IP strategy][Exists] check:  %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][IP strategy][Exists] check: id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("IP strategy id is borken, id is %s", rs.Primary.ID)
		}
		serviceId := idSplit[0]
		strategyId := idSplit[1]
		has, err := service.DescribeIPStrategyHas(ctx, serviceId, strategyId)
		if err != nil {
			return err
		}

		if !has {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][IP strategy][Exists] check: not exists: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAPIGatewayServiceBase = `
resource "tencentcloud_api_gateway_service" "service" {
  	service_name = "niceservice"
  	protocol     = "http&https"
  	service_desc = "your nice service"
  	net_type     = ["INNER", "OUTER"]
  	ip_version   = "IPv4"
}
`

const testApiIPStrategy_basic = testAPIGatewayServiceBase + `
resource "tencentcloud_api_gateway_ip_strategy" "test"{
    service_id    = tencentcloud_api_gateway_service.service.id 
    strategy_name = "tf_test"
    strategy_type = "BLACK"
    strategy_data = "9.9.9.9"
}
`

const testApiIPStrategy_update = testAPIGatewayServiceBase + `
resource "tencentcloud_api_gateway_ip_strategy" "test"{
    service_id    = tencentcloud_api_gateway_service.service.id 
    strategy_name = "tf_test"
    strategy_type = "BLACK"
    strategy_data = "7.7.7.7"
}
`
