package bh_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixDasbUserGroupResource_basic -v
func TestAccTencentCloudNeedFixDasbUserGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDasbUserGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_user_group.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_user_group.example", "name", "tf_example"),
				),
			},
			{
				ResourceName:      "tencentcloud_dasb_user_group.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDasbUserGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_user_group.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_user_group.example", "name", "tf_example_update"),
				),
			},
		},
	})
}

const testAccDasbUserGroup = `
resource "tencentcloud_dasb_user_group" "example" {
  name          = "tf_example"
}
`

const testAccDasbUserGroupUpdate = `
resource "tencentcloud_dasb_user_group" "example" {
  name          = "tf_example_update"
}
`
