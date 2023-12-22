package tco_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudOrganizationOrgMemberPolicyAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_ORGANIZATION)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgMemberPolicyAttachment,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_org_member_policy_attachment.org_member_policy_attachment", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_org_member_policy_attachment.org_member_policy_attachment", "member_uins.#"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_member_policy_attachment.org_member_policy_attachment", "policy_name", "example-iac"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_member_policy_attachment.org_member_policy_attachment", "identity_id", "1"),
				),
			},
		},
	})
}

const testAccOrganizationOrgMemberPolicyAttachment = `

resource "tencentcloud_organization_org_member_policy_attachment" "org_member_policy_attachment" {
  member_uins = [100033905366,100033905356]
  policy_name = "example-iac"
  identity_id = 1
}
`
