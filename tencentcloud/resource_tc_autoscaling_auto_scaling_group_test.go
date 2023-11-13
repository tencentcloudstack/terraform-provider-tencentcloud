package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudAutoscalingAutoScalingGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAutoscalingAutoScalingGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_autoscaling_auto_scaling_group.auto_scaling_group", "id")),
			},
			{
				ResourceName:      "tencentcloud_autoscaling_auto_scaling_group.auto_scaling_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccAutoscalingAutoScalingGroup = `

resource "tencentcloud_autoscaling_auto_scaling_group" "auto_scaling_group" {
  auto_scaling_group_id = "asg-xxxxxxxx"
  tags = {
    "createdBy" = "terraform"
  }
}

`
