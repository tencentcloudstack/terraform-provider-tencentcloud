package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudApigatewayUsagePlanResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApigatewayUsagePlan,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_apigateway_usage_plan.usage_plan", "id")),
			},
			{
				ResourceName:      "tencentcloud_apigateway_usage_plan.usage_plan",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccApigatewayUsagePlan = `

resource "tencentcloud_apigateway_usage_plan" "usage_plan" {
  usage_plan_name = ""
  usage_plan_desc = ""
  max_request_num = 
  max_request_num_pre_sec = 
}

`
