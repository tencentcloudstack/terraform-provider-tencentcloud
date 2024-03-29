package bh_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixDasbResetUserResource_basic -v
func TestAccTencentCloudNeedFixDasbResetUserResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDasbResetUser,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_reset_user.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_reset_user.example", "user_id"),
				),
			},
		},
	})
}

const testAccDasbResetUser = `
resource "tencentcloud_dasb_reset_user" "example" {
  user_id = 16
}
`
