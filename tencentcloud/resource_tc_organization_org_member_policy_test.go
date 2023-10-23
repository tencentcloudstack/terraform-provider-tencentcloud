package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudOrganizationOrgMemberPolicyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgMemberPolicy,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_org_member_policy.org_member_policy", "id")),
			},
			{
				ResourceName:      "tencentcloud_organization_org_member_policy.org_member_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccOrganizationOrgMemberPolicy = `

resource "tencentcloud_organization_org_member_policy" "org_member_policy" {
  member_uins = &lt;nil&gt;
  policy_name = &lt;nil&gt;
  identity_id = &lt;nil&gt;
  description = &lt;nil&gt;
}

`
