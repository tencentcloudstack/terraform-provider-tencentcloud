package as_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixAsScaleInInstancesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAsScaleInInstances,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_as_scale_in_instances.scale_in_instances", "id")),
			},
			{
				ResourceName:      "tencentcloud_as_scale_in_instances.scale_in_instances",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccAsScaleInInstances = `

resource "tencentcloud_as_scale_in_instances" "scale_in_instances" {
  auto_scaling_group_id = "asg-519acdug"
  scale_in_number = 1
}

`
