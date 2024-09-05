package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudIdentityCenterRoleConfigurationPermissionPolicyAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityCenterRoleConfigurationPermissionPolicyAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_role_configuration_permission_policy_attachment.identity_center_role_configuration_permission_policy_attachment", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_role_configuration_permission_policy_attachment.identity_center_role_configuration_permission_policy_attachment", "role_policy_name"),
				),
			},
			{
				ResourceName:      "tencentcloud_identity_center_role_configuration_permission_policy_attachment.identity_center_role_configuration_permission_policy_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccIdentityCenterRoleConfigurationPermissionPolicyAttachment = `
resource "tencentcloud_identity_center_role_configuration" "identity_center_role_configuration" {
    zone_id = "z-s64jh54hbcra"
    role_configuration_name = "tf-test-attachment"
    description = "test"
}

resource "tencentcloud_identity_center_role_configuration_permission_policy_attachment" "identity_center_role_configuration_permission_policy_attachment" {
    zone_id = "z-s64jh54hbcra"
    role_configuration_id = tencentcloud_identity_center_role_configuration.identity_center_role_configuration.role_configuration_id
    role_policy_id = 186457
	role_policy_name = "QcloudVPCReadOnlyAccess"
}
`
