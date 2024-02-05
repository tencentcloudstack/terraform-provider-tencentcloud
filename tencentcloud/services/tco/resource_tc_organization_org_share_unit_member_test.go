package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudOrganizationOrgShareUnitMemberResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgShareUnitMember,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_org_share_unit_member.org_share_unit_member", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_org_share_unit_member.org_share_unit_member", "unit_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_org_share_unit_member.org_share_unit_member", "area"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_org_share_unit_member.org_share_unit_member", "members.#"),
				),
			},
		},
	})
}

const testAccOrganizationOrgShareUnitMember = `
resource "tencentcloud_organization_org_share_unit" "org_share_unit" {
  name = "iac-test"
  area = "ap-guangzhou"
  description = "iac-test"
}
resource "tencentcloud_organization_org_share_unit_member" "org_share_unit_member" {
  unit_id = tencentcloud_organization_org_share_unit.org_share_unit.unit_id
  area = tencentcloud_organization_org_share_unit.org_share_unit.area
  members {
		share_member_uin=100035309479	
  }
}

`
