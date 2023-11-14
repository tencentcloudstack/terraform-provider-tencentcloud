package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudAutoscalingScaleInInstancesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAutoscalingScaleInInstances,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_autoscaling_scale_in_instances.scale_in_instances", "id")),
			},
			{
				ResourceName:      "tencentcloud_autoscaling_scale_in_instances.scale_in_instances",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccAutoscalingScaleInInstances = `

resource "tencentcloud_autoscaling_scale_in_instances" "scale_in_instances" {
  auto_scaling_group_id = "asg-xxxxxxxx"
  scale_in_number = 1
}

`
