package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcdbResetAccountPasswordResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbResetAccountPassword,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dcdb_reset_account_password.reset_account_password", "id")),
			},
			{
				ResourceName:      "tencentcloud_dcdb_reset_account_password.reset_account_password",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcdbResetAccountPassword = `

resource "tencentcloud_dcdb_reset_account_password" "reset_account_password" {
  instance_id = ""
  user_name = ""
  host = ""
  password = ""
}

`
