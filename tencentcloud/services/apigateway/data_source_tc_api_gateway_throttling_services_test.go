package apigateway_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testAPIGatewaythrottlingServiceDataSourceName = "data.tencentcloud_api_gateway_throttling_services"

func TestAccTencentAPIGatewayThrottlingServicesDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckAPIGatewayServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTestAccTencentAPIGatewayThrottlingServices(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAPIGatewayServiceExists("tencentcloud_api_gateway_service.service"),
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

func testAccTestAccTencentAPIGatewayThrottlingServices() string {
	return `
resource "tencentcloud_api_gateway_service" "service" {
  	service_name     = "niceservice"
  	protocol         = "http&https"
  	service_desc     = "your nice service"
  	net_type         = ["INNER", "OUTER"]
	ip_version       = "IPv4"
	release_limit    = 100
	pre_limit        = 100
	test_limit       = 100
}

data "tencentcloud_api_gateway_throttling_services" "id" {
    service_id = tencentcloud_api_gateway_service.service.id
}
`
}
