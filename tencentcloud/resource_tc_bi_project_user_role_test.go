package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudBiProjectUserRoleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBiProjectUserRole,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_bi_project_user_role.project_user_role", "id")),
			},
			{
				ResourceName:      "tencentcloud_bi_project_user_role.project_user_role",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccBiProjectUserRole = `

resource "tencentcloud_bi_project_user_role" "project_user_role" {
  project_id = 123
  user_list {
		user_id = "abc"
		user_name = "abc"
		corp_id = "abc"
		email = "abc@tencent.com"
		last_login = "2023-05-11 16:59:16"
		status = 1
		first_modify = 1
		phone_number = "12345678910"
		area_code = "86"
		created_user = "abc"
		created_at = "2023-05-11 16:59:16"
		updated_user = "abc"
		updated_at = "2023-05-11 16:59:16"
		global_user_name = "abc"
		mobile = "12345678910"

  }
  user_info_list {
		user_id = "abc"
		user_name = "abc"
		email = "abc@tencent.com"
		phone_number = "12345678910"
		area_code = "86"

  }
}

`
