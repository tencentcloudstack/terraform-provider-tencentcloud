package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudOrganizationOrgIdentityResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgIdentity,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_org_identity.org_identity", "id"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_identity.org_identity", "identity_alias_name", "example-iac-test"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_org_identity.org_identity", "identity_policy.#"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_identity.org_identity", "identity_policy.0.policy_id", "1"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_identity.org_identity", "identity_policy.0.policy_name", "AdministratorAccess"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_identity.org_identity", "identity_policy.0.policy_type", "2"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_identity.org_identity", "description", "iac-test"),
				),
			},
			{
				Config: testAccOrganizationOrgIdentityUpdate,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_org_identity.org_identity", "id"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_identity.org_identity", "identity_alias_name", "example-iac-test"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_org_identity.org_identity", "identity_policy.#"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_identity.org_identity", "identity_policy.0.policy_id", "2"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_identity.org_identity", "identity_policy.0.policy_name", "QCloudResourceFullAccess"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_identity.org_identity", "identity_policy.0.policy_type", "2"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_identity.org_identity", "description", "iac"),
				),
			},
			{
				ResourceName:      "tencentcloud_organization_org_identity.org_identity",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccOrganizationOrgIdentity = `

resource "tencentcloud_organization_org_identity" "org_identity" {
  identity_alias_name = "example-iac-test"
  identity_policy {
    policy_id = 1
    policy_name = "AdministratorAccess"
    policy_type = 2
  }
  description = "iac-test"
}

`
const testAccOrganizationOrgIdentityUpdate = `

resource "tencentcloud_organization_org_identity" "org_identity" {
  identity_alias_name = "example-iac-test"
  identity_policy {
    policy_id = 2
    policy_name = "QCloudResourceFullAccess"
    policy_type = 2
  }
  description = "iac"
}

`
