package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testAPIGatewayServicesDataSourceName = "data.tencentcloud_api_gateway_services"

func TestAccTencentAPIGatewayServicesDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAPIGatewayServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTestAccTencentAPIGatewayServices(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAPIGatewayServiceExists(testAPIGatewayServiceResourceName+".service"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServicesDataSourceName+".name", "list.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServicesDataSourceName+".name", "list.0.service_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServicesDataSourceName+".name", "list.0.service_name"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServicesDataSourceName+".name", "list.0.service_desc"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServicesDataSourceName+".name", "list.0.create_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServicesDataSourceName+".name", "list.0.inner_http_port"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServicesDataSourceName+".name", "list.0.inner_https_port"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServicesDataSourceName+".name", "list.0.internal_sub_domain"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServicesDataSourceName+".name", "list.0.ip_version"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServicesDataSourceName+".name", "list.0.net_type.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServicesDataSourceName+".name", "list.0.outer_sub_domain"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServicesDataSourceName+".name", "list.0.protocol"),

					resource.TestCheckResourceAttr(testAPIGatewayServicesDataSourceName+".id", "list.#", "1"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServicesDataSourceName+".id", "list.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServicesDataSourceName+".id", "list.0.service_id"),
					resource.TestCheckResourceAttr(testAPIGatewayServicesDataSourceName+".id", "list.0.service_name", "niceservice"),
					resource.TestCheckResourceAttr(testAPIGatewayServicesDataSourceName+".id", "list.0.service_desc", "your nice service"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServicesDataSourceName+".id", "list.0.create_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServicesDataSourceName+".id", "list.0.inner_http_port"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServicesDataSourceName+".id", "list.0.inner_https_port"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServicesDataSourceName+".id", "list.0.internal_sub_domain"),
					resource.TestCheckResourceAttr(testAPIGatewayServicesDataSourceName+".id", "list.0.ip_version", "IPv4"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServicesDataSourceName+".id", "list.0.net_type.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServicesDataSourceName+".id", "list.0.outer_sub_domain"),
					resource.TestCheckResourceAttr(testAPIGatewayServicesDataSourceName+".id", "list.0.protocol", "http&https"),
				),
			},
		},
	})
}

func testAccTestAccTencentAPIGatewayServices() string {
	return `
resource "tencentcloud_api_gateway_service" "service" {
  service_name = "niceservice"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

data "tencentcloud_api_gateway_services" "name" {
    service_name = tencentcloud_api_gateway_service.service.service_name
}

data "tencentcloud_api_gateway_services" "id" {
    service_id = tencentcloud_api_gateway_service.service.id
}
`
}
