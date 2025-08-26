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
		Steps: []resource.TestStep{{
			Config: testAccOrganizationMemberAuthPolicyAttachment,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_member_auth_policy_attachment.organization_member_auth_policy_attachment", "id")),
		}, {
			ResourceName:      "tencentcloud_organization_member_auth_policy_attachment.organization_member_auth_policy_attachment",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccOrganizationMemberAuthPolicyAttachment = `

resource "tencentcloud_organization_member_auth_policy_attachment" "organization_member_auth_policy_attachment" {
}
`
