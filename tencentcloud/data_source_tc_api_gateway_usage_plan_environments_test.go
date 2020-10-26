package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var (
	testAPIGatewayUsagePlansDataEnvironmentsSourceName = "data.tencentcloud_api_gateway_usage_plan_environments"
	testAPIGatewayUsagePlanAttachmentResource          = "tencentcloud_api_gateway_usage_plan_attachment.attach_service"
)

func TestAccTencentAPIGatewayUsagePlanEnvironmentsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAPIGatewayUsagePlanAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTestAccTencentAPIGatewayUsagePlanEnvironments(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAPIGatewayUsagePlanAttachmentExists(testAPIGatewayUsagePlanAttachmentResource),
					resource.TestCheckResourceAttrSet(testAPIGatewayUsagePlansDataEnvironmentsSourceName+".environment_test", "list.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewayUsagePlansDataEnvironmentsSourceName+".environment_test", "list.0.service_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewayUsagePlansDataEnvironmentsSourceName+".environment_test", "list.0.service_name"),
					resource.TestCheckResourceAttrSet(testAPIGatewayUsagePlansDataEnvironmentsSourceName+".environment_test", "list.0.environment"),
				),
			},
		},
	})
}

func testAccTestAccTencentAPIGatewayUsagePlanEnvironments() string {
	return `
		resource "tencentcloud_api_gateway_usage_plan" "plan" {
  			usage_plan_name         = "my_plan"
  			usage_plan_desc         = "nice plan"
  			max_request_num         = 100
  			max_request_num_pre_sec = 10
		}
		
		resource "tencentcloud_api_gateway_service" "service" {
			service_name = "niceservice"
			protocol     = "http&https"
			service_desc = "your nice service"
			net_type     = ["INNER", "OUTER"]
			ip_version   = "IPv4"
		}

		resource "tencentcloud_api_gateway_usage_plan_attachment" "attach_service" {
  			usage_plan_id  = tencentcloud_api_gateway_usage_plan.plan.id
  			service_id     = tencentcloud_api_gateway_service.service.id 
  			environment    = "test"
  			bind_type      = "SERVICE"
		}

		data "tencentcloud_api_gateway_usage_plan_environments" "environment_test" {
			usage_plan_id = tencentcloud_api_gateway_usage_plan_attachment.attach_service.usage_plan_id
			bind_type     = "SERVICE"
		}
	`
}
