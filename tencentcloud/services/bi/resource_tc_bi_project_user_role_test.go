package bi_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudBiProjectUserRoleResource_basic -v
func TestAccTencentCloudBiProjectUserRoleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBiProjectUserRole,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bi_project_user_role.project_user_role", "id"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project_user_role.project_user_role", "area_code", "+86"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project_user_role.project_user_role", "project_id", "11015030"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project_user_role.project_user_role", "role_id_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project_user_role.project_user_role", "email", "123456@qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project_user_role.project_user_role", "phone_number", "13130001000"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project_user_role.project_user_role", "user_id", "100024664626"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project_user_role.project_user_role", "user_name", "keep-cam-user"),
				),
			},
			{
				ResourceName:      "tencentcloud_bi_project_user_role.project_user_role",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBiProjectUserRoleUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bi_project_user_role.project_user_role", "id"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project_user_role.project_user_role", "area_code", "+86"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project_user_role.project_user_role", "project_id", "11015030"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project_user_role.project_user_role", "role_id_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project_user_role.project_user_role", "email", "1234567@qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project_user_role.project_user_role", "phone_number", "13130001000"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project_user_role.project_user_role", "user_id", "100024664626"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project_user_role.project_user_role", "user_name", "keep-cam-user"),
				),
			},
			{
				Config: testAccBiProjectUserRole,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bi_project_user_role.project_user_role", "id"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project_user_role.project_user_role", "area_code", "+86"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project_user_role.project_user_role", "project_id", "11015030"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project_user_role.project_user_role", "role_id_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project_user_role.project_user_role", "email", "123456@qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project_user_role.project_user_role", "phone_number", "13130001000"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project_user_role.project_user_role", "user_id", "100024664626"),
					resource.TestCheckResourceAttr("tencentcloud_bi_project_user_role.project_user_role", "user_name", "keep-cam-user"),
				),
			},
		},
	})
}

const testAccBiProjectUserRole = `

resource "tencentcloud_bi_project_user_role" "project_user_role" {
  area_code    = "+86"
  project_id   = 11015030
  role_id_list = [10629453]
  email        = "123456@qq.com"
  phone_number = "13130001000"
  user_id      = "100024664626"
  user_name    = "keep-cam-user"
}

`

const testAccBiProjectUserRoleUp = `

resource "tencentcloud_bi_project_user_role" "project_user_role" {
  area_code    = "+86"
  project_id   = 11015030
  role_id_list = [10629455]
  email        = "1234567@qq.com"
  phone_number = "13130001000"
  user_id      = "100024664626"
  user_name    = "keep-cam-user"
}

`
