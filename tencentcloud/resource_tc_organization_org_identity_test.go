package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
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
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_org_identity.org_identity", "id")),
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
  identity_alias_name = &lt;nil&gt;
  identity_policy {
		policy_id = &lt;nil&gt;
		policy_name = &lt;nil&gt;
		policy_type = &lt;nil&gt;
		policy_document = &lt;nil&gt;

  }
  description = &lt;nil&gt;
}

`
