package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudOrganizationOrgShareUnitNodeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgShareUnitNode,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_organization_org_share_unit_node.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_org_share_unit_node.example", "unit_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_org_share_unit_node.example", "node_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_organization_org_share_unit_node.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccOrganizationOrgShareUnitNode = `
resource "tencentcloud_organization_org_share_unit" "example" {
  name        = "tf-example"
  area        = "ap-guangzhou"
  description = "description."
}

resource "tencentcloud_organization_org_share_unit_node" "example" {
  unit_id = tencentcloud_organization_org_share_unit.example.unit_id
  node_id = 1001  # Replace with your actual node ID
}
`
