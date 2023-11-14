package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

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
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_bi_user_role.user_role", "id")),
			},
			{
				ResourceName:      "tencentcloud_bi_user_role.user_role",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccBiUserRole = `

resource "tencentcloud_bi_user_role" "user_role" {
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
