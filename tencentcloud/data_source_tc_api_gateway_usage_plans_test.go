package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testAPIGatewayUsagePlansDataSourceName = "data.tencentcloud_api_gateway_usage_plans"

func TestAccTencentAPIGatewayUsagePlansDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAPIGatewayUsagePlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTestAccTencentAPIGatewayUsagePlans(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAPIGatewayUsagePlanExists(testAPIGatewayUsagePlanResourceName+".plan"),
					resource.TestCheckResourceAttrSet(testAPIGatewayUsagePlansDataSourceName+".name", "list.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewayUsagePlansDataSourceName+".name", "list.0.usage_plan_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewayUsagePlansDataSourceName+".name", "list.0.usage_plan_name"),
					resource.TestCheckResourceAttrSet(testAPIGatewayUsagePlansDataSourceName+".name", "list.0.usage_plan_desc"),
					resource.TestCheckResourceAttrSet(testAPIGatewayUsagePlansDataSourceName+".name", "list.0.max_request_num"),
					resource.TestCheckResourceAttrSet(testAPIGatewayUsagePlansDataSourceName+".name", "list.0.max_request_num_pre_sec"),
					resource.TestCheckResourceAttrSet(testAPIGatewayUsagePlansDataSourceName+".name", "list.0.create_time"),
					resource.TestCheckResourceAttr(testAPIGatewayUsagePlansDataSourceName+".id", "list.#", "1"),
					resource.TestCheckResourceAttrSet(testAPIGatewayUsagePlansDataSourceName+".id", "list.0.usage_plan_id"),
					resource.TestCheckResourceAttr(testAPIGatewayUsagePlansDataSourceName+".id", "list.0.usage_plan_name", "my_plan"),
					resource.TestCheckResourceAttr(testAPIGatewayUsagePlansDataSourceName+".id", "list.0.usage_plan_desc", "nice plan"),
					resource.TestCheckResourceAttr(testAPIGatewayUsagePlansDataSourceName+".id", "list.0.max_request_num", "100"),
					resource.TestCheckResourceAttr(testAPIGatewayUsagePlansDataSourceName+".id", "list.0.max_request_num_pre_sec", "10"),
					resource.TestCheckResourceAttrSet(testAPIGatewayUsagePlansDataSourceName+".id", "list.0.create_time"),
				),
			},
		},
	})
}

func testAccTestAccTencentAPIGatewayUsagePlans() string {
	return `
		resource "tencentcloud_api_gateway_usage_plan" "plan" {
  			usage_plan_name         = "my_plan"
  			usage_plan_desc         = "nice plan"
  			max_request_num         = 100
  			max_request_num_pre_sec = 10
		}

		data "tencentcloud_api_gateway_usage_plans" "name" {
  			usage_plan_name = tencentcloud_api_gateway_usage_plan.plan.usage_plan_name
		}

		data "tencentcloud_api_gateway_usage_plans" "id" {
  			usage_plan_id = tencentcloud_api_gateway_usage_plan.plan.id
		}
	`
}
