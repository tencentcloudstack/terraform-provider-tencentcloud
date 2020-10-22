package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAPIGatewaythrottlingApiDataSourceName = "data.tencentcloud_api_gateway_throttling_apis"

func TestAccTencentAPIGatewayThrottlingApisDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckThrottlingAPIDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTestAccTencentAPIGatewayThrottlingApis(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckThrottlingAPIExists("tencentcloud_api_gateway_throttling_api.service"),
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

func testAccCheckThrottlingAPIDestroy(s *terraform.State) error {
	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		throttlingService = APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_api_gateway_throttling_api" {
			continue
		}

		serviceId := rs.Primary.Attributes["service_id"]
		environmentName := rs.Primary.Attributes["environment_name"]
		apiIds := rs.Primary.Attributes["api_ids"]
		environmentList, err := throttlingService.DescribeApiEnvironmentStrategyList(ctx, serviceId, []string{environmentName}, "")
		if err != nil {
			return err
		}

		for _, v := range environmentList {
			if v == nil || !strings.Contains(apiIds, *v.ApiId) {
				continue
			}
			environmentSet := v.EnvironmentStrategySet
			for _, env := range environmentSet {
				if env == nil || *env.EnvironmentName != environmentName {
					continue
				}

				if *env.Quota == QUOTA || *env.Quota == QUOTA_MAX {
					continue
				}
				return fmt.Errorf("throttling API still not restore: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckThrottlingAPIExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var (
			logId             = getLogId(contextNil)
			ctx               = context.WithValue(context.TODO(), logIdKey, logId)
			throttlingService = APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("API Getway throttling API %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("API Getway throttling API id is not set")
		}
		serviceId := rs.Primary.Attributes["service_id"]
		environmentName := rs.Primary.Attributes["environment_name"]
		apiIds := rs.Primary.Attributes["api_ids"]
		environmentList, err := throttlingService.DescribeApiEnvironmentStrategyList(ctx, serviceId, []string{environmentName}, "")
		if err != nil {
			return err
		}

		for _, v := range environmentList {
			if v == nil || !strings.Contains(apiIds, *v.ApiId) {
				continue
			}
			environmentSet := v.EnvironmentStrategySet
			for _, env := range environmentSet {
				if env == nil || *env.EnvironmentName != environmentName {
					continue
				}

				if *env.Quota == QUOTA {
					return fmt.Errorf("throttling API still not set value: %s", rs.Primary.ID)
				}
			}
		}
		return nil
	}
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
}

resource "tencentcloud_api_gateway_throttling_api" "service" {
	service_id       = tencentcloud_api_gateway_service.service.id
	strategy         = "400"
	environment_name = "test"
	api_ids          = [tencentcloud_api_gateway_api.api.id]
}

data "tencentcloud_api_gateway_throttling_apis" "id" {
    service_id = tencentcloud_api_gateway_throttling_api.service.service_id
}

data "tencentcloud_api_gateway_throttling_apis" "foo" {
	service_id        = tencentcloud_api_gateway_throttling_api.service.service_id
	environment_names = ["release", "test"]
}
`
}
