package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudProvisionRoleConfigurationOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProvisionRoleConfigurationOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_provision_role_configuration_operation.provision_role_configuration_operation", "id"),
				),
			},
		},
	})
}

const testAccProvisionRoleConfigurationOperation = `
resource "tencentcloud_identity_center_user" "identity_center_user" {
  zone_id     = "z-s64jh54hbcra"
  user_name   = "provision-role-configuration-user"
  description = "test"
}

resource "tencentcloud_identity_center_role_configuration" "identity_center_role_configuration" {
  zone_id                 = "z-s64jh54hbcra"
  role_configuration_name = "provision-role-configuration"
  description             = "test"
}

data "tencentcloud_organization_members" "members" {}

resource "tencentcloud_identity_center_role_assignment" "identity_center_role_assignment" {
  zone_id               = "z-s64jh54hbcra"
  principal_id          = tencentcloud_identity_center_user.identity_center_user.user_id
  principal_type        = "User"
  target_uin            = data.tencentcloud_organization_members.members.items.0.member_uin
  target_type           = "MemberUin"
  role_configuration_id = tencentcloud_identity_center_role_configuration.identity_center_role_configuration.role_configuration_id
}

resource "tencentcloud_identity_center_role_configuration_permission_policy_attachment" "identity_center_role_configuration_permission_policy_attachment" {
  zone_id               = "z-s64jh54hbcra"
  role_configuration_id = tencentcloud_identity_center_role_configuration.identity_center_role_configuration.role_configuration_id
  role_policy_id        = 186457
  role_policy_name      = "QcloudVPCReadOnlyAccess"
  depends_on            = [tencentcloud_identity_center_role_assignment.identity_center_role_assignment]
}
resource "tencentcloud_provision_role_configuration_operation" "provision_role_configuration_operation" {
  zone_id               = "z-s64jh54hbcra"
  role_configuration_id = tencentcloud_identity_center_role_configuration.identity_center_role_configuration.role_configuration_id
  target_type           = "MemberUin"
  target_uin            = "100038074533"
  depends_on            = [tencentcloud_identity_center_role_configuration_permission_policy_attachment.identity_center_role_configuration_permission_policy_attachment]
}
`
