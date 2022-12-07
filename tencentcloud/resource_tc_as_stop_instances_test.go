package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestTencentCloudAsStopInstancesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAsStopInstances,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_as_stop_instances.stop_instances", "id")),
			},
		},
	})
}

const testAccAsStopInstances = `

resource "tencentcloud_as_stop_instances" "stop_instances" {
  auto_scaling_group_id = ""
  instance_ids = ""
  stopped_mode = ""
}

`
