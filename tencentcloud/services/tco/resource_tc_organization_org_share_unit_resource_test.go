package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudOrganizationOrgShareUnitResourceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgShareUnitResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_organization_org_share_unit_resource.organization_org_share_unit_resource", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_organization_org_share_unit_resource.organization_org_share_unit_resource",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccOrganizationOrgShareUnitResource = `
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
`
