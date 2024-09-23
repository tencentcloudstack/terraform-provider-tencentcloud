package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudInviteOrganizationMemberOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInviteOrganizationMemberOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_invite_organization_member_operation.invite_organization_member_operation", "id"),
					resource.TestCheckResourceAttr("tencentcloud_invite_organization_member_operation.invite_organization_member_operation", "member_uin", "100038691413"),
					resource.TestCheckResourceAttr("tencentcloud_invite_organization_member_operation.invite_organization_member_operation", "name", "tf-test"),
					resource.TestCheckResourceAttr("tencentcloud_invite_organization_member_operation.invite_organization_member_operation", "policy_type", "Financial"),
					resource.TestCheckResourceAttr("tencentcloud_invite_organization_member_operation.invite_organization_member_operation", "node_id", "2002416"),
					resource.TestCheckResourceAttr("tencentcloud_invite_organization_member_operation.invite_organization_member_operation", "is_allow_quit", "Allow"),
					resource.TestCheckResourceAttr("tencentcloud_invite_organization_member_operation.invite_organization_member_operation", "permission_ids.#", "3"),
				),
			},
		},
	})
}

const testAccInviteOrganizationMemberOperation = `
resource "tencentcloud_invite_organization_member_operation" "invite_organization_member_operation" {
  member_uin = "100038691413"
  name = "tf-test"
  policy_type = "Financial"
  node_id = "2002416"
  is_allow_quit = "Allow"
  permission_ids = ["1", "2", "4"]
}
`
