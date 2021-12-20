package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testAPIGatewayAPIsDataSourceName = "data.tencentcloud_api_gateway_apis"

func TestAccTencentAPIGatewayAPIsDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAPIGatewayAPIDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTestAccTencentAPIGatewayAPIs(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAPIGatewayAPIExists(testAPIGatewayAPIResourceName+".api"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIsDataSourceName+".name", "list.0.api_name"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIsDataSourceName+".name", "list.0.api_desc"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIsDataSourceName+".name", "list.0.auth_type"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIsDataSourceName+".name", "list.0.protocol"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIsDataSourceName+".name", "list.0.enable_cors"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIsDataSourceName+".name", "list.0.request_config_path"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIsDataSourceName+".name", "list.0.request_config_method"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIsDataSourceName+".name", "list.0.service_config_type"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIsDataSourceName+".name", "list.0.service_config_timeout"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIsDataSourceName+".name", "list.0.service_config_url"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIsDataSourceName+".name", "list.0.service_config_path"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIsDataSourceName+".name", "list.0.service_config_method"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIsDataSourceName+".name", "list.0.response_type"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIsDataSourceName+".name", "list.0.response_success_example"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIsDataSourceName+".name", "list.0.response_fail_example"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIsDataSourceName+".name", "list.0.modify_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIsDataSourceName+".name", "list.0.create_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIsDataSourceName+".name", "list.0.create_time"),

					resource.TestCheckResourceAttr(testAPIGatewayAPIsDataSourceName+".id", "list.#", "1"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIsDataSourceName+".id", "list.0.api_name", "hello"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIsDataSourceName+".id", "list.0.api_desc", "my hello api"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIsDataSourceName+".id", "list.0.auth_type", "NONE"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIsDataSourceName+".id", "list.0.protocol", "HTTP"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIsDataSourceName+".id", "list.0.enable_cors", "true"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIsDataSourceName+".id", "list.0.request_config_path", "/user/info"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIsDataSourceName+".id", "list.0.request_config_method", "GET"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIsDataSourceName+".id", "list.0.service_config_type", "HTTP"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIsDataSourceName+".id", "list.0.service_config_timeout", "15"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIsDataSourceName+".id", "list.0.service_config_url", "http://www.qq.com"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIsDataSourceName+".id", "list.0.service_config_path", "/user"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIsDataSourceName+".id", "list.0.service_config_method", "GET"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIsDataSourceName+".id", "list.0.response_type", "HTML"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIsDataSourceName+".id", "list.0.response_success_example", "success"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIsDataSourceName+".id", "list.0.response_fail_example", "fail"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIsDataSourceName+".id", "list.0.modify_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIsDataSourceName+".id", "list.0.create_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIsDataSourceName+".id", "list.0.create_time"),
				),
			},
		},
	})
}

func testAccTestAccTencentAPIGatewayAPIs() string {
	return `
resource "tencentcloud_api_gateway_service" "service" {
  service_name = "ck"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_api" "api" {
  service_id               = tencentcloud_api_gateway_service.service.id
  api_name                 = "hello"
  api_desc                 = "my hello api"
  auth_type                = "NONE"
  protocol                 = "HTTP"
  enable_cors              = true
  request_config_path      = "/user/info"
  request_config_method    = "GET"
  service_config_type      = "HTTP"
  service_config_timeout   = 15
  service_config_url       = "http://www.qq.com"
  service_config_path      = "/user"
  service_config_method    = "GET"
  response_type            = "HTML"
  response_success_example = "success"
  response_fail_example    = "fail"
}

data "tencentcloud_api_gateway_apis" "id" {
  service_id = tencentcloud_api_gateway_service.service.id
  api_id     = tencentcloud_api_gateway_api.api.id
}

data "tencentcloud_api_gateway_apis" "name" {
  service_id = tencentcloud_api_gateway_service.service.id
  api_name   = tencentcloud_api_gateway_api.api.api_name
}
`
}
