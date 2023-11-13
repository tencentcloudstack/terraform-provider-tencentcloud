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
					resource.TestCheckResourceAttr("tencentcloud_dasb_user.example", "user_name", "tf_example"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_user.example", "real_name", "terraform"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_user.example", "phone", "+86|18345678782"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_user.example", "email", "demo@tencent.com"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_user.example", "validate_from", "2023-09-22T02:00:00+08:00"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_user.example", "validate_to", "2023-09-23T03:00:00+08:00"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_user.example", "department_id", "1.2"),
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
					resource.TestCheckResourceAttr("tencentcloud_dasb_user.example", "user_name", "tf_example"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_user.example", "real_name", "terraform_update"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_user.example", "phone", "+86|1834567880"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_user.example", "email", "demo1@tencent.com"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_user.example", "validate_from", "2023-09-22T01:00:00+08:00"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_user.example", "validate_to", "2023-09-23T02:00:00+08:00"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_user.example", "department_id", "1.2"),
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
  phone         = "+86|18345678782"
  email         = "demo@tencent.com"
  validate_from = "2023-09-22T02:00:00+08:00"
  validate_to   = "2023-09-23T03:00:00+08:00"
  department_id = "1.2"
  auth_type     = 0
}
`

const testAccDasbUserUpdate = `
resource "tencentcloud_dasb_user" "example" {
  user_name     = "tf_example"
  real_name     = "terraform_update"
  phone         = "+86|1834567880"
  email         = "demo1@tencent.com"
  validate_from = "2023-09-22T01:00:00+08:00"
  validate_to   = "2023-09-23T02:00:00+08:00"
  department_id = "1.2"
  auth_type     = 0
}
`
