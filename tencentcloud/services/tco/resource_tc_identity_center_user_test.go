package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudIdentityCenterUserResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityCenterUser,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_user.identity_center_user", "id"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_user.identity_center_user", "zone_id", "z-s64jh54hbcra"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_user.identity_center_user", "user_name", "tf-test-user"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_user.identity_center_user", "description", "test"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_user.identity_center_user", "user_status"),
				),
			},
			{
				Config: testAccIdentityCenterUserUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_user.identity_center_user", "id"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_user.identity_center_user", "description", "test-update"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_user.identity_center_user", "user_status", "Disabled"),
				),
			},
			{
				ResourceName:      "tencentcloud_identity_center_user.identity_center_user",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccIdentityCenterUser = `
resource "tencentcloud_identity_center_user" "identity_center_user" {
    zone_id = "z-s64jh54hbcra"
    user_name = "tf-test-user"
    description = "test"
}
`

const testAccIdentityCenterUserUpdate = `
resource "tencentcloud_identity_center_user" "identity_center_user" {
    zone_id = "z-s64jh54hbcra"
    user_name = "tf-test-user"
    description = "test-update"
	user_status = "Disabled"
}
`
