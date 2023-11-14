package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudAutoscalingLastActivityDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAutoscalingLastActivityDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_autoscaling_last_activity.last_activity")),
			},
		},
	})
}

const testAccAutoscalingLastActivityDataSource = `

data "tencentcloud_autoscaling_last_activity" "last_activity" {
  auto_scaling_group_ids = &lt;nil&gt;
  }

`
