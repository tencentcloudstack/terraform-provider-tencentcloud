package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudIdentityCenterRoleAssignmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityCenterRoleAssignment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_role_assignment.identity_center_role_assignment", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_role_assignment.identity_center_role_assignment", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_role_assignment.identity_center_role_assignment", "principal_id"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_role_assignment.identity_center_role_assignment", "principal_type", "User"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_role_assignment.identity_center_role_assignment", "role_configuration_id"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_role_assignment.identity_center_role_assignment", "target_type", "MemberUin"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_role_assignment.identity_center_role_assignment", "target_uin"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_role_assignment.identity_center_role_assignment", "update_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_role_assignment.identity_center_role_assignment", "zone_id"),
				),
			},
			{
				ResourceName: "tencentcloud_identity_center_role_assignment.identity_center_role_assignment",
				ImportState:  true,
			},
		},
	})
}

const testAccIdentityCenterRoleAssignment = `
resource "tencentcloud_identity_center_user" "identity_center_user" {
    zone_id = "z-s64jh54hbcra"
    user_name = "assignment-user"
    description = "test"
}

resource "tencentcloud_identity_center_role_configuration" "identity_center_role_configuration" {
    zone_id = "z-s64jh54hbcra"
    role_configuration_name = "assignment-configuration"
    description = "test"
}

data "tencentcloud_organization_members" "members" {}

resource "tencentcloud_identity_center_role_assignment" "identity_center_role_assignment" {
  zone_id = "z-s64jh54hbcra"
  principal_id = tencentcloud_identity_center_user.identity_center_user.user_id
  principal_type = "User"
  target_uin = data.tencentcloud_organization_members.members.items.0.member_uin
  target_type = "MemberUin"
  role_configuration_id = tencentcloud_identity_center_role_configuration.identity_center_role_configuration.role_configuration_id
}
`
