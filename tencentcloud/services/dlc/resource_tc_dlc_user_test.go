package dlc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcUserResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcUser,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_user.user", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_user.user", "user_description", "for terraform test"),
				),
			},
			{
				Config: testAccDlcUserUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_user.user", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_user.user", "user_description", "for terraform"),
				),
			},
			{
				ResourceName:      "tencentcloud_dlc_user.user",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDlcUser = `

resource "tencentcloud_dlc_user" "user" {
  user_id          = "100027012454"
  user_type        = "COMMON"
  user_alias       = "terraform-test"
  user_description = "for terraform test"
}

`

const testAccDlcUserUpdate = `

resource "tencentcloud_dlc_user" "user" {
  user_id          = "100027012454"
  user_type        = "COMMON"
  user_alias       = "terraform-test"
  user_description = "for terraform"
}

`
