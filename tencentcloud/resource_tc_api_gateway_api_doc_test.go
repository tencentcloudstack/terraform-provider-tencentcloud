package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

var testAPIGatewayAPIDocResourceName = "tencentcloud_api_gateway_api_doc"
var testAPIGatewayAPIDocResourceKey = testAPIGatewayAPIDocResourceName + ".test"

func TestAccTencentCloudAPIGateWayAPIDocResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAPIGatewayAPIDocDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGatewayAPIDoc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGatewayAPIDocExists(testAPIGatewayAPIDocResourceKey),
					resource.TestCheckResourceAttr(testAPIGatewayAPIDocResourceKey, "api_doc_name", "test"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIDocResourceKey, "service_id", "service-2nuhovb7"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIDocResourceKey, "environment", "release"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIDocResourceKey, "api_ids", ""),
				),
			},
			{
				ResourceName:      testAPIGatewayAPIDocResourceKey,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAPIGatewayAPIDocDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAPIGatewayAPIDocResourceName {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		apiDoc, err := service.DescribeApiDoc(ctx, rs.Primary.ID)
		if apiDoc != nil {
			return fmt.Errorf("api_gateway api_doc %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckAPIGatewayAPIDocExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		command, err := service.DescribeApiDoc(ctx, rs.Primary.ID)
		if command == nil {
			return fmt.Errorf("api_gateway api_doc %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccAPIGatewayAPIDoc = `
resource "tencentcloud_api_gateway_api_doc" "test" {
  api_doc_name = "doc_test1"
  service_id   = "service_test1"
  environment  = "release"
  api_ids      = ["api-test1", "api-test2"]
}
`
