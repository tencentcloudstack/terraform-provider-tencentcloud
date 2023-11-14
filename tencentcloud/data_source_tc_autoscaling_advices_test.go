package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudAutoscalingAdvicesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAutoscalingAdvicesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_autoscaling_advices.advices")),
			},
		},
	})
}

const testAccAutoscalingAdvicesDataSource = `

data "tencentcloud_autoscaling_advices" "advices" {
  auto_scaling_group_ids = &lt;nil&gt;
  }

`
