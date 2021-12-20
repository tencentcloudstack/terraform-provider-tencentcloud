package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var (
	testAPIGatewayUsagePlanResourceName = "tencentcloud_api_gateway_usage_plan"
	testAPIGatewayUsagePlanResourceKey  = testAPIGatewayUsagePlanResourceName + ".plan"
)

func TestAccTencentCloudAPIGateWayUsagePlanResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAPIGatewayUsagePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGatewayUsagePlan,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGatewayUsagePlanExists(testAPIGatewayUsagePlanResourceKey),
					resource.TestCheckResourceAttr(testAPIGatewayUsagePlanResourceKey, "usage_plan_name", "my_plan"),
					resource.TestCheckResourceAttr(testAPIGatewayUsagePlanResourceKey, "usage_plan_desc", "nice plan"),
					resource.TestCheckResourceAttr(testAPIGatewayUsagePlanResourceKey, "max_request_num", "100"),
					resource.TestCheckResourceAttr(testAPIGatewayUsagePlanResourceKey, "max_request_num_pre_sec", "10"),
					resource.TestCheckResourceAttrSet(testAPIGatewayUsagePlanResourceKey, "create_time"),
				),
			},
			{
				ResourceName:      testAPIGatewayUsagePlanResourceKey,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAPIGatewayUsagePlanUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGatewayUsagePlanExists(testAPIGatewayUsagePlanResourceKey),
					resource.TestCheckResourceAttr(testAPIGatewayUsagePlanResourceKey, "usage_plan_name", "my_plan_update"),
					resource.TestCheckResourceAttr(testAPIGatewayUsagePlanResourceKey, "usage_plan_desc", "nice plan update"),
					resource.TestCheckResourceAttr(testAPIGatewayUsagePlanResourceKey, "max_request_num", "10"),
					resource.TestCheckResourceAttr(testAPIGatewayUsagePlanResourceKey, "max_request_num_pre_sec", "100"),
					resource.TestCheckResourceAttrSet(testAPIGatewayUsagePlanResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayUsagePlanResourceKey, "modify_time"),
				),
			},
		},
	})
}

func testAccCheckAPIGatewayUsagePlanDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAPIGatewayUsagePlanResourceName {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeUsagePlan(ctx, rs.Primary.ID)
		if err != nil {
			_, has, err = service.DescribeUsagePlan(ctx, rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		if !has {
			return nil
		} else {
			return fmt.Errorf("delete API gateway usage plan %s fail, still on server", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckAPIGatewayUsagePlanExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeUsagePlan(ctx, rs.Primary.ID)
		if err != nil {
			_, has, err = service.DescribeUsagePlan(ctx, rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		if has {
			return nil
		} else {
			return fmt.Errorf("API gateway usage plan %s not found on server", rs.Primary.ID)

		}
	}
}

const testAccAPIGatewayUsagePlan = `
	resource "tencentcloud_api_gateway_usage_plan" "plan" {
  		usage_plan_name         = "my_plan"
  		usage_plan_desc         = "nice plan"
  		max_request_num         = 100
  		max_request_num_pre_sec = 10
	}
`
const testAccAPIGatewayUsagePlanUpdate = `
	resource "tencentcloud_api_gateway_usage_plan" "plan" {
		usage_plan_name         = "my_plan_update"
		usage_plan_desc         = "nice plan update"
		max_request_num         = 10
		max_request_num_pre_sec = 100
	}
`
