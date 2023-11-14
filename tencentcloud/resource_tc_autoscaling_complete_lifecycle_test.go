package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudAutoscalingCompleteLifecycleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAutoscalingCompleteLifecycle,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_autoscaling_complete_lifecycle.complete_lifecycle", "id")),
			},
			{
				ResourceName:      "tencentcloud_autoscaling_complete_lifecycle.complete_lifecycle",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccAutoscalingCompleteLifecycle = `

resource "tencentcloud_autoscaling_complete_lifecycle" "complete_lifecycle" {
  lifecycle_hook_id = "ash-xxxxxxxx"
  lifecycle_action_result = "CONTINUE"
  instance_id = "ins-xxxxxxxx"
  lifecycle_action_token = &lt;nil&gt;
}

`
