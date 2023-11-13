package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudOrganizationOrgMemberEmailResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgMemberEmail,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_org_member_email.org_member_email", "id")),
			},
			{
				ResourceName:      "tencentcloud_organization_org_member_email.org_member_email",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccOrganizationOrgMemberEmail = `

resource "tencentcloud_organization_org_member_email" "org_member_email" {
  member_uin = &lt;nil&gt;
  email = &lt;nil&gt;
  country_code = &lt;nil&gt;
  phone = &lt;nil&gt;
  bind_id = &lt;nil&gt;
}

`
