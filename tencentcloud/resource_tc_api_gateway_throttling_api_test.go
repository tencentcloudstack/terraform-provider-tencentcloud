package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudAPIGateWayThrottlingAPI(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckThrottlingAPIDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccThrottlingAPI,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckThrottlingAPIExists("tencentcloud_api_gateway_throttling_api.foo"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_throttling_api.foo", "service_id"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_throttling_api.foo", "strategy", "400"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_throttling_api.foo", "environment_name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_throttling_api.foo", "api_ids.#", "1"),
				),
			},
			{
				Config: testAccThrottlingAPIUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckThrottlingAPIExists("tencentcloud_api_gateway_throttling_api.foo"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_throttling_api.foo", "service_id"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_throttling_api.foo", "strategy", "400"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_throttling_api.foo", "environment_name", "release"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_throttling_api.foo", "api_ids.#", "1"),
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
		environmentList, err := throttlingService.DescribeApiEnvironmentStrategyList(ctx, serviceId, []string{environmentName})
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
		environmentList, err := throttlingService.DescribeApiEnvironmentStrategyList(ctx, serviceId, []string{environmentName})
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

const testAPIBase = `
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
`

const testAccThrottlingAPI = testAPIBase + `
resource "tencentcloud_api_gateway_throttling_api" "foo" {
	service_id       = tencentcloud_api_gateway_service.service.id 
	strategy         = "400"
	environment_name = "test"
	api_ids          = [tencentcloud_api_gateway_api.api.id]
}
`

const testAccThrottlingAPIUpdate = testAPIBase + `
resource "tencentcloud_api_gateway_throttling_api" "foo" {
	service_id       = tencentcloud_api_gateway_service.service.id
	strategy         = "400"
	environment_name = "release"
	api_ids          = [tencentcloud_api_gateway_api.api.id]
}
`
