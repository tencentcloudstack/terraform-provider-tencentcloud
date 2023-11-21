package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDasbResetUserResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDasbResetUser,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_reset_user.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_reset_user.example", "user_id"),
				),
			},
		},
	})
}

const testAccDasbResetUser = `
resource "tencentcloud_dasb_reset_user" "example" {
  user_id = 16
}
`
