package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCfwBlockIgnoreListResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwBlockIgnoreList,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cfw_block_ignore_list.block_ignore_list", "id")),
			},
			{
				ResourceName:      "tencentcloud_cfw_block_ignore_list.block_ignore_list",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCfwBlockIgnoreList = `

resource "tencentcloud_cfw_block_ignore_list" "block_ignore_list" {
  rules {
		direction = 0
		end_time = "2023-09-09 15:04:05"
		i_p = "1.1.1.1"
		domain = ""
		comment = "block ip 1.1.1.1"
		start_time = "2023-09-01 15:04:05"

  }
  rule_type = 1
}

`
