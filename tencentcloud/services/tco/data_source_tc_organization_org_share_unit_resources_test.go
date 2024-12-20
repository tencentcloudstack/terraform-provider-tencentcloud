package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudOrganizationOrgShareUnitResourcesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgShareUnitResourcesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_organization_org_share_unit_resources.organization_org_share_unit_resources"),
					resource.TestCheckResourceAttr("data.tencentcloud_organization_org_share_unit_resources.organization_org_share_unit_resources", "items.#", "1"),
				),
			},
		},
	})
}

const testAccOrganizationOrgShareUnitResourcesDataSource = `
resource "tencentcloud_organization_org_share_unit" "org_share_unit" {
  name = "iac-test"
  area = "ap-guangzhou"
  description = "iac-test"
}

resource "tencentcloud_organization_org_share_unit_resource" "organization_org_share_unit_resource" {
  unit_id = split("#", tencentcloud_organization_org_share_unit.org_share_unit.id)[1]
  area = "ap-guangzhou"
  type = "secret"
  product_resource_id = "100027395662/keep-tf"
}

data "tencentcloud_organization_org_share_unit_resources" "organization_org_share_unit_resources" {
  area = "ap-guangzhou"
  unit_id = split("#", tencentcloud_organization_org_share_unit.org_share_unit.id)[1]
  depends_on = [ tencentcloud_organization_org_share_unit_resource.organization_org_share_unit_resource ]
}
`
