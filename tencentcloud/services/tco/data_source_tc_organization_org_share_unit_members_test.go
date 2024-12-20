package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudOrganizationOrgShareUnitMembersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgShareUnitMembersDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_organization_org_share_unit_members.organization_org_share_unit_members"),
					resource.TestCheckResourceAttr("data.tencentcloud_organization_org_share_unit_members.organization_org_share_unit_members", "items.#", "1"),
				),
			},
		},
	})
}

const testAccOrganizationOrgShareUnitMembersDataSource = `
resource "tencentcloud_organization_org_share_unit" "org_share_unit" {
  name = "iac-test"
  area = "ap-guangzhou"
  description = "iac-test"
}

resource "tencentcloud_organization_org_share_unit_member" "org_share_unit_member" {
  unit_id = tencentcloud_organization_org_share_unit.org_share_unit.unit_id
  area = tencentcloud_organization_org_share_unit.org_share_unit.area
  members {
	share_member_uin=100038074517	
  }
}

data "tencentcloud_organization_org_share_unit_members" "organization_org_share_unit_members" {
  unit_id = split("#", tencentcloud_organization_org_share_unit.org_share_unit.id)[1]
  area = "ap-guangzhou"
  search_key = "100038074517"
  depends_on = [ tencentcloud_organization_org_share_unit_member.org_share_unit_member ]
}
`
