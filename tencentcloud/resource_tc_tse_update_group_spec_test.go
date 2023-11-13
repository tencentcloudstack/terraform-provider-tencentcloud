package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTseUpdateGroupSpecResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseUpdateGroupSpec,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tse_update_group_spec.update_group_spec", "id")),
			},
			{
				ResourceName:      "tencentcloud_tse_update_group_spec.update_group_spec",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTseUpdateGroupSpec = `

resource "tencentcloud_tse_update_group_spec" "update_group_spec" {
  gateway_id = ""
  group_id = ""
  node_config {
		specification = ""
		number = 

  }
}

`
