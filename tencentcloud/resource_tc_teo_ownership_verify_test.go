package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTeoOwnershipVerifyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoOwnershipVerify,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_ownership_verify.ownership_verify", "id")),
			},
			{
				ResourceName:      "tencentcloud_teo_ownership_verify.ownership_verify",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoOwnershipVerify = `

resource "tencentcloud_teo_ownership_verify" "ownership_verify" {
  domain = "qq.com"
}

`
