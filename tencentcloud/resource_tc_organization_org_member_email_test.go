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
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_org_member_email.org_member_email", "id"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_member_email.org_member_email", "member_uin", "100033704327"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_member_email.org_member_email", "email", "iac-test@qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_member_email.org_member_email", "country_code", "86"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_member_email.org_member_email", "phone", "12345678901"),
				),
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
  member_uin = 100033704327
  email = "iac-test@qq.com"
  country_code = "86"
  phone = "12345678901"
  }

`
