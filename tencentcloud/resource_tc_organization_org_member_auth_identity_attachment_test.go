package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudOrganizationOrgMemberAuthIdentityResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgMemberAuthIdentity,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_org_member_auth_identity_attachment.org_member_auth_identity", "id"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_member_auth_identity_attachment.org_member_auth_identity", "member_uin", "100033704327"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_org_member_auth_identity_attachment.org_member_auth_identity", "identity_ids.#"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_member_auth_identity_attachment.org_member_auth_identity", "identity_ids.0", "1657"),
				),
			},
			{
				ResourceName:      "tencentcloud_organization_org_member_auth_identity_attachment.org_member_auth_identity",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccOrganizationOrgMemberAuthIdentity = `

resource "tencentcloud_organization_org_member_auth_identity_attachment" "org_member_auth_identity" {
  member_uin = 100033704327
  identity_ids = [1657]
}

`
