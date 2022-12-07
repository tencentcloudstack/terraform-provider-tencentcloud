package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestTencentCloudAsStartInstancesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAsStartInstances,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_as_start_instances.start_instances", "id")),
			},
		},
	})
}

const testAccAsStartInstances = `

resource "tencentcloud_as_start_instances" "start_instances" {
  auto_scaling_group_id = ""
  instance_ids = ""
}

`
