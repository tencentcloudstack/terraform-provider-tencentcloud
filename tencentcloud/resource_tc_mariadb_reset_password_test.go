package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMariadbResetPasswordResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbResetPassword,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mariadb_reset_password.reset_password", "id")),
			},
			{
				ResourceName:      "tencentcloud_mariadb_reset_password.reset_password",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbResetPassword = `

resource "tencentcloud_mariadb_reset_password" "reset_password" {
  instance_id = ""
  user_name = ""
  host = ""
  password = ""
}

`
