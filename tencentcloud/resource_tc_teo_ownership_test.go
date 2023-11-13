package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTeoOwnershipResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoOwnership,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_ownership.ownership", "id")),
			},
			{
				ResourceName:      "tencentcloud_teo_ownership.ownership",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoOwnership = `

resource "tencentcloud_teo_ownership" "ownership" {
  domain = "qq.com"
}

`
