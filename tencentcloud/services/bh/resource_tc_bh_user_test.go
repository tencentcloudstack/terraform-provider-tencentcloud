package bh_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudBhUserResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBhUser,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_bh_user.example", "id")),
			},
			{
				Config: testAccBhUserUpdate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_bh_user.example", "id")),
			},
			{
				ResourceName:      "tencentcloud_bh_user.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccBhUser = `
resource "tencentcloud_bh_user" "example" {
  user_name = "tf-example"
  real_name = "Terraform"
  phone     = "+86|18991162528"
  email     = "demo@tencent.com"
}
`

const testAccBhUserUpdate = `
resource "tencentcloud_bh_user" "example" {
  user_name = "tf-example"
  real_name = "Terraform-updated"
  phone     = "+86|18991162528"
  email     = "demo@tencent.com"
}
`
