package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataProjectMemberResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataProjectMember,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_project_member.example", "id"),
				),
			},
			{
				Config: testAccWedataProjectMemberUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_project_member.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_project_member.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWedataProjectMember = `
resource "tencentcloud_wedata_project_member" "example" {
  project_id = "2983848457986924544"
  user_uin   = "100044238258"
  role_ids = [
    "308335260274237440",
    "308335260844662784"
  ]
}
`

const testAccWedataProjectMemberUpdate = `
resource "tencentcloud_wedata_project_member" "example" {
  project_id = "2983848457986924544"
  user_uin   = "100044238258"
  role_ids = [
    "308335260274237440",
  ]
}
`
