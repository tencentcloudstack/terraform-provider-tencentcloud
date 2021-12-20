package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testAPIGatewaythrottlingApiDataSourceName = "data.tencentcloud_api_gateway_throttling_apis"

func TestAccTencentAPIGatewayThrottlingApisDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAPIGatewayAPIDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTestAccTencentAPIGatewayThrottlingApis(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAPIGatewayAPIExists("tencentcloud_api_gateway_api.api"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".id", "list.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".id", "list.0.service_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".id", "list.0.api_environment_strategies.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".id", "list.0.api_environment_strategies.0.api_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".id", "list.0.api_environment_strategies.0.api_name"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".id", "list.0.api_environment_strategies.0.path"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".id", "list.0.api_environment_strategies.0.method"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".id", "list.0.api_environment_strategies.0.strategy_list.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".id", "list.0.api_environment_strategies.0.strategy_list.0.environment_name"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".id", "list.0.api_environment_strategies.0.strategy_list.0.quota"),

					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".foo", "list.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".foo", "list.0.service_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".foo", "list.0.api_environment_strategies.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".foo", "list.0.api_environment_strategies.0.api_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".foo", "list.0.api_environment_strategies.0.api_name"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".foo", "list.0.api_environment_strategies.0.path"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".foo", "list.0.api_environment_strategies.0.method"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".foo", "list.0.api_environment_strategies.0.strategy_list.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".foo", "list.0.api_environment_strategies.0.strategy_list.0.environment_name"),
					resource.TestCheckResourceAttrSet(testAPIGatewaythrottlingApiDataSourceName+".foo", "list.0.api_environment_strategies.0.strategy_list.0.quota"),
				),
			},
		},
	})
}

func testAccTestAccTencentAPIGatewayThrottlingApis() string {
	return `
resource "tencentcloud_api_gateway_service" "service" {
  	service_name = "niceservice"
  	protocol     = "http&https"
  	service_desc = "your nice service"
  	net_type     = ["INNER", "OUTER"]
  	ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_api" "api" {
    service_id            = tencentcloud_api_gateway_service.service.id
    api_name              = "hello_update"
    api_desc              = "my hello api update"
    auth_type             = "SECRET"
    protocol              = "HTTP"
    enable_cors           = true
    request_config_path   = "/user/info"
    request_config_method = "POST"
    request_parameters {
    	name          = "email"
        position      = "QUERY"
        type          = "string"
        desc          = "your email please?"
        default_value = "tom@qq.com"
        required      = true
    }
    service_config_type      = "HTTP"
    service_config_timeout   = 10
    service_config_url       = "http://www.tencent.com"
    service_config_path      = "/user"
    service_config_method    = "POST"
    response_type            = "XML"
    response_success_example = "<note>success</note>"
	response_fail_example    = "<note>fail</note>"
	
    response_error_codes {
    	code           = 10
        msg            = "system error"
       	desc           = "system error code"
       	converted_code = -10
        need_convert   = true
	}
	
	release_limit    = 100
	pre_limit        = 100
	test_limit       = 100
}

data "tencentcloud_api_gateway_throttling_apis" "id" {
    service_id = tencentcloud_api_gateway_api.api.service_id
}

data "tencentcloud_api_gateway_throttling_apis" "foo" {
	service_id        = tencentcloud_api_gateway_api.api.service_id
	environment_names = ["release", "test"]
}
`
}
