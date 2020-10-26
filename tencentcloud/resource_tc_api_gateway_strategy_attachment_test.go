package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudAPIGateWayStrategyAttachment_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testApiStrategyAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testApiStrategyAttachment_basic,
				Check: resource.ComposeTestCheckFunc(
					testApiStrategyAttachmentExists("tencentcloud_api_gateway_strategy_attachment.test"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_strategy_attachment.test", "service_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_strategy_attachment.test", "strategy_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_strategy_attachment.test", "environment_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_strategy_attachment.test", "bind_api_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_api_gateway_strategy_attachment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testApiStrategyAttachmentDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_api_gateway_strategy_attachment" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 4 {
			return fmt.Errorf("IP strategy attachment id is broken, id is %s", rs.Primary.ID)
		}
		serviceId := idSplit[0]
		strategyId := idSplit[1]
		bindApiId := idSplit[2]

		has, err := service.DescribeStrategyAttachment(ctx, serviceId, strategyId, bindApiId)
		if err != nil {
			return err
		}

		if has {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][IP strategy][Destroy] check: IP strategy still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testApiStrategyAttachmentExists(n string) resource.TestCheckFunc {
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
		if len(idSplit) != 4 {
			return fmt.Errorf("IP strategy attachment id is broken, id is %s", rs.Primary.ID)
		}
		serviceId := idSplit[0]
		strategyId := idSplit[1]
		bindApiId := idSplit[2]
		has, err := service.DescribeStrategyAttachment(ctx, serviceId, strategyId, bindApiId)
		if err != nil {
			return err
		}

		if !has {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][IP strategy][Exists] check: not exists: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAPIGatewayServiceAttachmentBase = `
resource "tencentcloud_api_gateway_service" "service" {
  	service_name = "niceservice"
  	protocol     = "http&https"
  	service_desc = "your nice service"
  	net_type     = ["INNER", "OUTER"]
  	ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_ip_strategy" "test"{
    service_id    = tencentcloud_api_gateway_service.service.id
    strategy_name = "tf_test"
    strategy_type = "BLACK"
    strategy_data = "9.9.9.9"
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

const testApiStrategyAttachment_basic = testAPIGatewayServiceAttachmentBase + `
resource "tencentcloud_api_gateway_strategy_attachment" "test"{
   service_id       = tencentcloud_api_gateway_service_release.service.service_id
   strategy_id      = tencentcloud_api_gateway_ip_strategy.test.strategy_id 
   environment_name = "release"
   bind_api_id      = tencentcloud_api_gateway_api.api.id
}
`
