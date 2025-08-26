package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudOrganizationMemberAuthPolicyAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationMemberAuthPolicyAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_organization_member_auth_policy_attachment.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_member_auth_policy_attachment.example", "policy_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_member_auth_policy_attachment.example", "org_sub_account_uin"),
				),
			},
			{
				ResourceName:      "tencentcloud_organization_member_auth_policy_attachment.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccOrganizationMemberAuthPolicyAttachment = `
resource "tencentcloud_organization_member_auth_policy_attachment" "example" {
  policy_id           = 250021751
  org_sub_account_uin = 100037718139
}
`
