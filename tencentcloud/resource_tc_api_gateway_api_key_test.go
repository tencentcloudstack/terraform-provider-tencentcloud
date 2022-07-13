package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAPIGatewayAPIKeyResourceName = "tencentcloud_api_gateway_api_key"
var testAPIGatewayAPIKeyResourceKey = testAPIGatewayAPIKeyResourceName + ".test"

func TestAccTencentCloudAPIGateWayAPIKeyResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAPIGatewayAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGatewayAPIKey,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGatewayAPIKeyExists(testAPIGatewayAPIKeyResourceKey),
					resource.TestCheckResourceAttr(testAPIGatewayAPIKeyResourceKey, "secret_name", "my_api_key"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIKeyResourceKey, "status", "on"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKey, "access_key_secret"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKey, "modify_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKey, "create_time"),
				),
			},
			{
				ResourceName:      testAPIGatewayAPIKeyResourceKey,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAPIGatewayAPIKeyUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGatewayAPIKeyExists(testAPIGatewayAPIKeyResourceKey),
					resource.TestCheckResourceAttr(testAPIGatewayAPIKeyResourceKey, "secret_name", "my_api_key"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIKeyResourceKey, "status", "off"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKey, "access_key_secret"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKey, "modify_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKey, "create_time"),
				),
			},
		},
	})
}

func testAccCheckAPIGatewayAPIKeyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAPIGatewayAPIKeyResourceName {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeApiKey(ctx, rs.Primary.ID)
		if err != nil {
			_, has, err = service.DescribeApiKey(ctx, rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		if !has {
			return nil
		} else {
			return fmt.Errorf("delete API key for API gateway %s fail, still on server", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckAPIGatewayAPIKeyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeApiKey(ctx, rs.Primary.ID)
		if err != nil {
			_, has, err = service.DescribeApiKey(ctx, rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		if has {
			return nil
		} else {
			return fmt.Errorf("API key for API gateway %s not found on server", rs.Primary.ID)

		}
	}
}

const testAccAPIGatewayAPIKey = `
resource "tencentcloud_api_gateway_api_key" "test" {
  secret_name = "my_api_key"
  status      = "on"
}
`
const testAccAPIGatewayAPIKeyUpdate = `
resource "tencentcloud_api_gateway_api_key" "test" {
  secret_name = "my_api_key"
  status      = "off"
}
`
