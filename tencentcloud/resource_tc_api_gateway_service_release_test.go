package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
)

var (
	testAPIGatewayServiceReleaseResourceName = "tencentcloud_api_gateway_service_release"
	testAPIGatewayServiceReleaseResourceKey  = testAPIGatewayServiceReleaseResourceName + ".service"
)

func TestAccTencentCloudAPIGatewayServiceRelease(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAPIGatewayServiceReleaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGatewayServiceRelease,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGatewayServiceReleaseExists(testAPIGatewayServiceReleaseResourceKey),
					resource.TestCheckResourceAttrSet(testAPIGatewayServiceReleaseResourceKey, "service_id"),
					resource.TestCheckResourceAttr(testAPIGatewayServiceReleaseResourceKey, "environment_name", "release"),
					resource.TestCheckResourceAttr(testAPIGatewayServiceReleaseResourceKey, "release_desc", "test service release"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServiceReleaseResourceKey, "release_version"),
				),
			},
			{
				ResourceName:      testAPIGatewayServiceReleaseResourceKey,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAPIGatewayServiceReleaseDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAPIGatewayServiceReleaseResourceName {
			continue
		}

		var (
			logId             = getLogId(contextNil)
			ctx               = context.WithValue(context.TODO(), logIdKey, logId)
			apiGatewayService = APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
			info              []*apigateway.ServiceReleaseHistoryInfo
			err               error
			has               bool
		)

		ids := strings.Split(rs.Primary.ID, FILED_SP)
		if len(ids) != 3 {
			return fmt.Errorf("id is broken, id is %s", rs.Primary.ID)
		}
		var (
			serviceId  = ids[0]
			envName    = ids[1]
			envVersion = ids[2]
		)

		info, has, err = apiGatewayService.DescribeServiceEnvironmentReleaseHistory(ctx, serviceId, envName)
		if err != nil {
			info, _, err = apiGatewayService.DescribeServiceEnvironmentReleaseHistory(ctx, serviceId, envName)
		}
		if err != nil {
			return err
		}
		if !has {
			return nil
		}

		for _, v := range info {
			if *v.VersionName == envVersion {
				return fmt.Errorf("API gateway service release %s fail, still on server", rs.Primary.ID)
			}
		}
		return nil
	}
	return nil
}

func testAccCheckAPIGatewayServiceReleaseExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}

		var (
			logId             = getLogId(contextNil)
			ctx               = context.WithValue(context.TODO(), logIdKey, logId)
			apiGatewayService = APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
			info              []*apigateway.ServiceReleaseHistoryInfo
			err               error
			has               bool
		)

		ids := strings.Split(rs.Primary.ID, FILED_SP)
		if len(ids) != 3 {
			return fmt.Errorf("id is broken, id is %s", rs.Primary.ID)
		}
		var (
			serviceId  = ids[0]
			envName    = ids[1]
			envVersion = ids[2]
		)

		info, has, err = apiGatewayService.DescribeServiceEnvironmentReleaseHistory(ctx, serviceId, envName)
		if err != nil {
			info, has, err = apiGatewayService.DescribeServiceEnvironmentReleaseHistory(ctx, serviceId, envName)
		}
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("API gateway service release %s not exist on server", rs.Primary.ID)
		}

		for _, v := range info {
			if *v.VersionName == envVersion {
				return nil
			}
		}

		return fmt.Errorf("API gateway service release %s not exist on server", rs.Primary.ID)
	}
}

const testAccAPIGatewayServiceRelease = `
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
    	code           = 20
        msg            = "system error"
       	desc           = "system error code"
       	converted_code = -10
        need_convert   = true
	}
}

resource "tencentcloud_api_gateway_service_release" "service" {
  service_id       = tencentcloud_api_gateway_api.api.service_id
  environment_name = "release"
  release_desc     = "test service release"
}
`
