package tco_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudOrganizationOrgMemberEmailResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
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
				Config: testAccOrganizationOrgMemberEmailUpdate,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_org_member_email.org_member_email", "id"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_member_email.org_member_email", "member_uin", "100033704327"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_member_email.org_member_email", "email", "iac-test-update@qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_member_email.org_member_email", "country_code", "86"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_member_email.org_member_email", "phone", "12345678902"),
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
const testAccOrganizationOrgMemberEmailUpdate = `

resource "tencentcloud_organization_org_member_email" "org_member_email" {
  member_uin = 100033704327
  email = "iac-test-update@qq.com"
  country_code = "86"
  phone = "12345678902"
  }

`
