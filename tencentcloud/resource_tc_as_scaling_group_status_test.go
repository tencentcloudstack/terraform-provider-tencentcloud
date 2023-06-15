package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAsScalingGroupStatusResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAsScalingGroupStatus,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_as_scaling_group_status.scaling_group_status", "id")),
			},
			{
				ResourceName:      "tencentcloud_as_scaling_group_status.scaling_group_status",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccAsScalingGroupStatus = `

resource "tencentcloud_as_scaling_group_status" "scaling_group_status" {
  auto_scaling_group_id = "asg-519acdug"
  enable                = false
}

`
