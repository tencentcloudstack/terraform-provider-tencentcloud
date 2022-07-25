package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAPIGatewayServiceResourceName = "tencentcloud_api_gateway_service"
var testAPIGatewayServiceResourceKey = testAPIGatewayServiceResourceName + ".service"

func TestAccTencentCloudNeedFixAPIGateWayServiceResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAPIGatewayServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGatewayService,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGatewayServiceExists(testAPIGatewayServiceResourceKey),
					resource.TestCheckResourceAttr(testAPIGatewayServiceResourceKey, "service_name", "myservice"),
					resource.TestCheckResourceAttr(testAPIGatewayServiceResourceKey, "protocol", "http"),
					resource.TestCheckResourceAttr(testAPIGatewayServiceResourceKey, "service_desc", "my nice service"),
					resource.TestCheckResourceAttr(testAPIGatewayServiceResourceKey, "ip_version", "IPv4"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServiceResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServiceResourceKey, "internal_sub_domain"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServiceResourceKey, "inner_http_port"),
					resource.TestCheckResourceAttr(testAPIGatewayServiceResourceKey, "release_limit", "500"),
					resource.TestCheckResourceAttr(testAPIGatewayServiceResourceKey, "pre_limit", "500"),
					resource.TestCheckResourceAttr(testAPIGatewayServiceResourceKey, "test_limit", "500"),
				),
			},
			{
				ResourceName:      testAPIGatewayServiceResourceKey,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAPIGatewayServiceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGatewayServiceExists(testAPIGatewayServiceResourceKey),
					resource.TestCheckResourceAttr(testAPIGatewayServiceResourceKey, "service_name", "yourservice"),
					resource.TestCheckResourceAttr(testAPIGatewayServiceResourceKey, "protocol", "http&https"),
					resource.TestCheckResourceAttr(testAPIGatewayServiceResourceKey, "service_desc", "your nice service"),
					resource.TestCheckResourceAttr(testAPIGatewayServiceResourceKey, "ip_version", "IPv4"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServiceResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServiceResourceKey, "modify_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServiceResourceKey, "internal_sub_domain"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServiceResourceKey, "outer_sub_domain"),
					resource.TestCheckResourceAttrSet(testAPIGatewayServiceResourceKey, "inner_http_port"),
					resource.TestCheckResourceAttr(testAPIGatewayServiceResourceKey, "release_limit", "100"),
					resource.TestCheckResourceAttr(testAPIGatewayServiceResourceKey, "pre_limit", "100"),
					resource.TestCheckResourceAttr(testAPIGatewayServiceResourceKey, "test_limit", "100"),
				),
			},
		},
	})
}

func testAccCheckAPIGatewayServiceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAPIGatewayServiceResourceName {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeService(ctx, rs.Primary.ID)
		if err != nil {
			_, has, err = service.DescribeService(ctx, rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		if !has {
			return nil
		} else {
			return fmt.Errorf("delete API gateway service %s fail, still on server", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckAPIGatewayServiceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeService(ctx, rs.Primary.ID)
		if err != nil {
			_, has, err = service.DescribeService(ctx, rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		if has {
			return nil
		} else {
			return fmt.Errorf("API gateway service %s not found on server", rs.Primary.ID)
		}
	}
}

const testAccAPIGatewayService = `
resource "tencentcloud_api_gateway_service" "service" {
  service_name     = "myservice"
  protocol         = "http"
  service_desc     = "my nice service"
  net_type         = ["INNER"]
  ip_version       = "IPv4"
  release_limit    = 500
  pre_limit        = 500
  test_limit       = 500
}
`
const testAccAPIGatewayServiceUpdate = `
resource "tencentcloud_api_gateway_service" "service" {
  service_name     = "yourservice"
  protocol         = "http&https"
  service_desc     = "your nice service"
  net_type         = ["INNER", "OUTER"]
  ip_version       = "IPv4"
  release_limit    = 100
  pre_limit        = 100
  test_limit       = 100
}
`
