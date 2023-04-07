package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainVerifyUserAccountResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainVerifyUserAccount,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_verify_user_account.verify_user_account", "id")),
			},
			{
				ResourceName:      "tencentcloud_dbbrain_verify_user_account.verify_user_account",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDbbrainVerifyUserAccount = `

resource "tencentcloud_dbbrain_verify_user_account" "verify_user_account" {
  instance_id = ""
  user = ""
  password = ""
  product = ""
}

`
