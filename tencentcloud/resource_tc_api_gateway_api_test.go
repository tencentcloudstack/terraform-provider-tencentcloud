package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAPIGatewayAPIResourceName = "tencentcloud_api_gateway_api"
var testAPIGatewayAPIResourceKey = testAPIGatewayAPIResourceName + ".api"

func TestAccTencentCloudAPIGateWayAPIResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAPIGatewayAPIDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGatewayBase,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGatewayAPIExists(testAPIGatewayAPIResourceKey),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "api_name", "hello"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "api_desc", "my hello api"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "auth_type", "NONE"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "enable_cors", "true"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "request_config_path", "/user/info"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "request_config_method", "GET"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "service_config_type", "HTTP"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "service_config_timeout", "15"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "service_config_url", "http://www.qq.com"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "service_config_path", "/user"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "service_config_method", "GET"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "response_type", "HTML"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "response_success_example", "success"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "response_fail_example", "fail"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIResourceKey, "modify_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIResourceKey, "create_time"),
				),
			},
			{
				Config: testAccAPIGatewayAPIUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGatewayAPIExists(testAPIGatewayAPIResourceKey),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "api_name", "hello_update"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "api_desc", "my hello api update"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "auth_type", "NONE"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "enable_cors", "true"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "request_config_path", "/user/info"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "request_config_method", "POST"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "service_config_type", "HTTP"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "service_config_timeout", "10"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "service_config_url", "http://www.tencent.com"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "service_config_path", "/user"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "service_config_method", "POST"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "response_type", "XML"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "response_success_example", "<note>success</note>"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIResourceKey, "response_fail_example", "<note>fail</note>"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIResourceKey, "modify_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIResourceKey, "create_time"),
				),
			},
		},
	})
}

func testAccCheckAPIGatewayAPIDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAPIGatewayAPIResourceName {
			continue
		}

		var (
			logId     = getLogId(contextNil)
			ctx       = context.WithValue(context.TODO(), logIdKey, logId)
			service   = APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
			apiId     = rs.Primary.ID
			serviceId = rs.Primary.Attributes["service_id"]
		)

		_, has, err := service.DescribeApi(ctx, serviceId, apiId)
		if err != nil {
			_, has, err = service.DescribeApi(ctx, serviceId, apiId)
		}
		if err != nil {
			return err
		}
		if !has {
			return nil
		} else {
			return fmt.Errorf("delete API for API gateway %s fail, still on server", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckAPIGatewayAPIExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}

		var (
			logId     = getLogId(contextNil)
			ctx       = context.WithValue(context.TODO(), logIdKey, logId)
			service   = APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
			apiId     = rs.Primary.ID
			serviceId = rs.Primary.Attributes["service_id"]
		)

		_, has, err := service.DescribeApi(ctx, serviceId, apiId)
		if err != nil {
			_, has, err = service.DescribeApi(ctx, serviceId, apiId)
		}
		if err != nil {
			return err
		}
		if has {
			return nil
		} else {
			return fmt.Errorf("API for API gateway %s not found on server", rs.Primary.ID)

		}
	}
}

const testAccAPIGatewayBase = `
resource "tencentcloud_api_gateway_service" "service" {
  service_name = "ck"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_api" "api" {
  service_id            = tencentcloud_api_gateway_service.service.id
  api_name              = "hello"
  api_desc              = "my hello api"
  auth_type             = "NONE"
  protocol              = "HTTP"
  enable_cors           = true
  request_config_path   = "/user/info"
  request_config_method = "GET"

  request_parameters {
    name          = "name"
    position      = "QUERY"
    type          = "string"
    desc          = "who are you?"
    default_value = "tom"
    required      = true
  }
  service_config_type      = "HTTP"
  service_config_timeout   = 15
  service_config_url       = "http://www.qq.com"
  service_config_path      = "/user"
  service_config_method    = "GET"
  response_type            = "HTML"
  response_success_example = "success"
  response_fail_example    = "fail"
  response_error_codes {
    code           = 100
    msg            = "system error"
    desc           = "system error code"
    converted_code = -100
    need_convert   = true
  }
}
`

const testAccAPIGatewayAPIUpdate = `
resource "tencentcloud_api_gateway_service" "service" {
  service_name = "ck"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_api" "api" {
  service_id            = tencentcloud_api_gateway_service.service.id
  api_name              = "hello_update"
  api_desc              = "my hello api update"
  auth_type             = "NONE"
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
