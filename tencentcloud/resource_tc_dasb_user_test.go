package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixDasbUserResource_basic -v
func TestAccTencentCloudNeedFixDasbUserResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDasbUser,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_user.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_user.example", "user_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_user.example", "real_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_user.example", "phone"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_user.example", "email"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_user.example", "validate_from"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_user.example", "validate_to"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_user.example", "auth_type"),
				),
			},
			{
				ResourceName:      "tencentcloud_dasb_user.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDasbUserUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_user.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_user.example", "user_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_user.example", "real_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_user.example", "phone"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_user.example", "email"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_user.example", "validate_from"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_user.example", "validate_to"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_user.example", "auth_type"),
				),
			},
		},
	})
}

const testAccDasbUser = `
resource "tencentcloud_dasb_user" "example" {
  user_name     = "tf_example"
  real_name     = "terraform"
  phone         = "+86|1834567878"
  email         = "demo@tencent.com"
  validate_from = "2023-09-22T00:00:00+08:00"
  validate_to   = "2023-09-23T00:00:00+08:00"
  auth_type     = 0
}
`

const testAccDasbUserUpdate = `
resource "tencentcloud_dasb_user" "example" {
  user_name     = "tf_example"
  real_name     = "terraform_update"
  phone         = "+86|1834567878"
  email         = "demo@tencent.com"
  validate_from = "2023-09-22T00:00:00+08:00"
  validate_to   = "2023-09-23T00:00:00+08:00"
  auth_type     = 0
}
`
