package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudBiUserRoleResource_basic -v
func TestAccTencentCloudBiUserRoleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBiUserRole,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bi_user_role.user_role", "id"),
					resource.TestCheckResourceAttr("tencentcloud_bi_user_role.user_role", "area_code", "+86"),
					resource.TestCheckResourceAttr("tencentcloud_bi_user_role.user_role", "email", "1055000000@qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_bi_user_role.user_role", "phone_number", "13470010000"),
					resource.TestCheckResourceAttr("tencentcloud_bi_user_role.user_role", "role_id_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_bi_user_role.user_role", "user_id", "100032767426"),
					resource.TestCheckResourceAttr("tencentcloud_bi_user_role.user_role", "user_name", "keep-iac-test"),
				),
			},
			{
				ResourceName:      "tencentcloud_bi_user_role.user_role",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBiUserRoleUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bi_user_role.user_role", "id"),
					resource.TestCheckResourceAttr("tencentcloud_bi_user_role.user_role", "area_code", "+86"),
					resource.TestCheckResourceAttr("tencentcloud_bi_user_role.user_role", "email", "1055000001@qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_bi_user_role.user_role", "phone_number", "13470010001"),
					resource.TestCheckResourceAttr("tencentcloud_bi_user_role.user_role", "role_id_list.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_bi_user_role.user_role", "user_id", "100032767426"),
					resource.TestCheckResourceAttr("tencentcloud_bi_user_role.user_role", "user_name", "keep-iac-test"),
				),
			},
		},
	})
}

const testAccBiUserRole = `

resource "tencentcloud_bi_user_role" "user_role" {
  area_code    = "+86"
  email        = "1055000000@qq.com"
  phone_number = "13470010000"
  role_id_list = [
    10629359,
  ]
  user_id   = "100032767426"
  user_name = "keep-iac-test"
}

`

const testAccBiUserRoleUp = `

resource "tencentcloud_bi_user_role" "user_role" {
  area_code    = "+86"
  email        = "1055000001@qq.com"
  phone_number = "13470010001"
  role_id_list = [
    10629360, 10629453,
  ]
  user_id   = "100032767426"
  user_name = "keep-iac-test"
}

`
