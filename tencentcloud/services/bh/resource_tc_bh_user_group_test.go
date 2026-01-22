package bh_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudBhUserGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBhUserGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bh_user_group.example", "id"),
				),
			},
			{
				Config: testAccBhUserGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bh_user_group.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_bh_user_group.example",
				ImportState:       true,
				ImportStateVerify: true,
			}},
	})
}

const testAccBhUserGroup = `
resource "tencentcloud_bh_user_group" "example" {
  name = "tf-example"
}
`

const testAccBhUserGroupUpdate = `
resource "tencentcloud_bh_user_group" "example" {
  name = "tf-example-updated"
}
`
