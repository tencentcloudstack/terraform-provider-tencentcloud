package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfReleaseApiGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfReleaseApiGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_release_api_group.release_api_group", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_release_api_group.release_api_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfReleaseApiGroup = `

resource "tencentcloud_tsf_release_api_group" "release_api_group" {
  group_id = "grp-qp0rj3zi"
}

`
