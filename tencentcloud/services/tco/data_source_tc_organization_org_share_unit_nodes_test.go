package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudOrganizationOrgShareUnitNodesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgShareUnitNodesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_organization_org_share_unit_nodes.organization_org_share_unit_nodes"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_organization_org_share_unit_nodes.organization_org_share_unit_nodes", "items.#"),
				),
			},
		},
	})
}

const testAccOrganizationOrgShareUnitNodesDataSource = `
resource "tencentcloud_organization_org_share_unit" "org_share_unit" {
  name        = "iac-test"
  area        = "ap-guangzhou"
  description = "iac-test"
}

resource "tencentcloud_organization_org_share_unit_node" "org_share_unit_node" {
  unit_id = tencentcloud_organization_org_share_unit.org_share_unit.unit_id
  node_id = 1001  # Replace with your actual node ID
}

data "tencentcloud_organization_org_share_unit_nodes" "organization_org_share_unit_nodes" {
  unit_id    = tencentcloud_organization_org_share_unit.org_share_unit.unit_id
  search_key = "1001"
  depends_on = [tencentcloud_organization_org_share_unit_node.org_share_unit_node]
}
`
