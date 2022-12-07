package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestTencentCloudAsProtectInstancesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAsProtectInstances,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_as_protect_instances.protect_instances", "id")),
			},
		},
	})
}

const testAccAsProtectInstances = `

resource "tencentcloud_as_protect_instances" "protect_instances" {
  auto_scaling_group_id = ""
  instance_ids = ""
  protected_from_scale_in = ""
}

`
