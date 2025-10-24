package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudIdentityCenterRoleConfigurationPermissionCustomPoliciesAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityCenterRoleConfigurationPermissionPoliciesAttachmentCustomPolicy,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_role_configuration_permission_custom_policies_attachment.identity_center_role_configuration_permission_custom_policies_attachment", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_role_configuration_permission_custom_policies_attachment.identity_center_role_configuration_permission_custom_policies_attachment", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_role_configuration_permission_custom_policies_attachment.identity_center_role_configuration_permission_custom_policies_attachment", "role_configuration_id"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_role_configuration_permission_custom_policies_attachment.identity_center_role_configuration_permission_custom_policies_attachment", "policies.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_identity_center_role_configuration_permission_custom_policies_attachment.identity_center_role_configuration_permission_custom_policies_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccIdentityCenterRoleConfigurationPermissionPoliciesAttachmentCustomPolicy = `
resource "tencentcloud_identity_center_role_configuration" "identity_center_role_configuration" {
    zone_id = "z-s64jh54hbcra"
    role_configuration_name = "tf-test-custom-policy"
    description = "test"
}

resource "tencentcloud_identity_center_role_configuration_permission_custom_policies_attachment" "identity_center_role_configuration_permission_custom_policies_attachment" {
    zone_id = "z-s64jh54hbcra"
    role_configuration_id = tencentcloud_identity_center_role_configuration.identity_center_role_configuration.role_configuration_id
    policies {
        role_policy_name = "CustomPolicy1"
        role_policy_document = <<-EOF
{
    "version": "2.0",
    "statement": [
        {
            "effect": "allow",
            "action": [
                "vpc:AcceptAttachCcnInstances"
            ],
            "resource": [
                "*"
            ]
        }
    ]
}
EOF
    }
    
}
`
