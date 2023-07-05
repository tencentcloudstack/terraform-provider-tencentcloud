package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCiamUserGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCiamUserGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ciam_user_group.user_group", "id")),
			},
			{
				ResourceName:      "tencentcloud_ciam_user_group.user_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCiamUserGroup = `

resource "tencentcloud_ciam_user_group" "user_group" {
  display_name = ""
  user_store_id = ""
  description = ""
}

`
