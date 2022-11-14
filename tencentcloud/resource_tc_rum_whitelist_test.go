package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudRumWhitelist_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumWhitelist,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_rum_whitelist.whitelist", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_rum_whitelist.whitelist",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRumWhitelist = `

resource "tencentcloud_rum_whitelist" "whitelist" {
  instance_i_d = ""
  remark = ""
  whitelist_uin = ""
  aid = ""
        }

`
