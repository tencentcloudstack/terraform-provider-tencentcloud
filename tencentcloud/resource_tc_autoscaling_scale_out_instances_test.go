package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudAutoscalingScaleOutInstancesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAutoscalingScaleOutInstances,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_autoscaling_scale_out_instances.scale_out_instances", "id")),
			},
			{
				ResourceName:      "tencentcloud_autoscaling_scale_out_instances.scale_out_instances",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccAutoscalingScaleOutInstances = `

resource "tencentcloud_autoscaling_scale_out_instances" "scale_out_instances" {
  auto_scaling_group_id = "asg-xxxxxxxx"
  scale_out_number = 1
}

`
