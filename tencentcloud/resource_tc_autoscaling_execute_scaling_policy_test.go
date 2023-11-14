package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudAutoscalingExecuteScalingPolicyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAutoscalingExecuteScalingPolicy,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_autoscaling_execute_scaling_policy.execute_scaling_policy", "id")),
			},
			{
				ResourceName:      "tencentcloud_autoscaling_execute_scaling_policy.execute_scaling_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccAutoscalingExecuteScalingPolicy = `

resource "tencentcloud_autoscaling_execute_scaling_policy" "execute_scaling_policy" {
  auto_scaling_policy_id = "asp-xxxxxxxx"
  honor_cooldown = false
  trigger_source = "API"
}

`
