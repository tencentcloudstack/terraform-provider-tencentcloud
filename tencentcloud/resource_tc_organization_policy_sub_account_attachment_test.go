package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudOrganizationPolicySubAccountAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationPolicySubAccountAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_policy_sub_account_attachment.policy_sub_account_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_organization_policy_sub_account_attachment.policy_sub_account_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccOrganizationPolicySubAccountAttachment = `

resource "tencentcloud_organization_policy_sub_account_attachment" "policy_sub_account_attachment" {
  policy_id = &lt;nil&gt;
  org_sub_account_uins = &lt;nil&gt;
  member_uin = &lt;nil&gt;
  org_sub_account_uin = &lt;nil&gt;
  policy_name = &lt;nil&gt;
  identity_id = &lt;nil&gt;
  identity_role_name = &lt;nil&gt;
  identity_role_alias_name = &lt;nil&gt;
  create_time = &lt;nil&gt;
  update_time = &lt;nil&gt;
  org_sub_account_name = &lt;nil&gt;
}

`
