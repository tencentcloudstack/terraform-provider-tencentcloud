package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var (
	testAPIGatewayIpStrategySourceName = "data.tencentcloud_api_gateway_ip_strategies"
	testAPIGatewayIpStrategy           = "tencentcloud_api_gateway_ip_strategy.test"
)

func TestAccTencentAPIGatewayIpStrategyDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testApiIPStrategyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTestAccTencentAPIGatewayIpStrategy(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testApiIPStrategyExists(testAPIGatewayIpStrategy),
					resource.TestCheckResourceAttrSet(testAPIGatewayIpStrategySourceName+".id", "list.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewayIpStrategySourceName+".id", "list.0.strategy_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewayIpStrategySourceName+".id", "list.0.strategy_name"),
					resource.TestCheckResourceAttrSet(testAPIGatewayIpStrategySourceName+".id", "list.0.bind_api_total_count"),
					resource.TestCheckResourceAttrSet(testAPIGatewayIpStrategySourceName+".id", "list.0.create_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayIpStrategySourceName+".id", "list.0.attach_list.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewayIpStrategySourceName+".name", "list.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewayIpStrategySourceName+".name", "list.0.strategy_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewayIpStrategySourceName+".name", "list.0.strategy_name"),
					resource.TestCheckResourceAttrSet(testAPIGatewayIpStrategySourceName+".name", "list.0.bind_api_total_count"),
					resource.TestCheckResourceAttrSet(testAPIGatewayIpStrategySourceName+".name", "list.0.create_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayIpStrategySourceName+".name", "list.0.attach_list.#"),
				),
			},
		},
	})
}

func testAccTestAccTencentAPIGatewayIpStrategy() string {
	return `
resource "tencentcloud_api_gateway_service" "service" {
  service_name = "ck"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_ip_strategy" "test"{
	service_id      = tencentcloud_api_gateway_service.service.id
	strategy_name	= "tf_test"
	strategy_type	= "BLACK"
	strategy_data	= "9.9.9.9"
}

data "tencentcloud_api_gateway_ip_strategies" "id" {
	service_id = tencentcloud_api_gateway_ip_strategy.test.service_id
}

data "tencentcloud_api_gateway_ip_strategies" "name" {
    service_id = tencentcloud_api_gateway_ip_strategy.test.service_id 
	strategy_name = tencentcloud_api_gateway_ip_strategy.test.strategy_name
}
`
}
