package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfStartContainerGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfStartContainerGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_start_container_group.start_container_group", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_start_container_group.start_container_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfStartContainerGroup = `

resource "tencentcloud_tsf_start_container_group" "start_container_group" {
  group_id = "group-xxxxxxxx"
}

`
