package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudOrganizationOrgShareUnitsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgShareUnitsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_organization_org_share_units.organization_org_share_units"),
					resource.TestCheckResourceAttr("data.tencentcloud_organization_org_share_units.organization_org_share_units", "items.#", "1"),
				),
			},
		},
	})
}

const testAccOrganizationOrgShareUnitsDataSource = `
resource "tencentcloud_organization_org_share_unit" "org_share_unit" {
  name = "iac-test-1"
  area = "ap-guangzhou"
  description = "iac-test"
}

data "tencentcloud_organization_org_share_units" "organization_org_share_units" {
  area = "ap-guangzhou"
  search_key = split("#", tencentcloud_organization_org_share_unit.org_share_unit.id)[1]
}
`
