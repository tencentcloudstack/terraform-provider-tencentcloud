package tco_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudOrganizationOrgNode_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgNode,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_organization_org_node.org_node", "id"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_node.org_node", "parent_node_id", "2014254"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_node.org_node", "name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_node.org_node", "remark", "for terraform test"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_node.org_node", "tags.tf-test1", "org_node1"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_node.org_node", "tags.tf-test2", "org_node2"),
				),
			},
			{
				Config: testAccOrganizationOrgNodeUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_organization_org_node.org_node", "id"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_node.org_node", "tags.tf-test1", "org_node_update"),
				),
			},
			{
				ResourceName:      "tencentcloud_organization_org_node.org_node",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccOrganizationOrgNode = `

resource "tencentcloud_organization_org_node" "org_node" {
  name           = "terraform_test"
  parent_node_id = 2014254
  remark         = "for terraform test"
    tags = {
    "tf-test1" = "org_node1"
	"tf-test2" = "org_node2"
  }
}

`

const testAccOrganizationOrgNodeUpdate = `

resource "tencentcloud_organization_org_node" "org_node" {
  name           = "terraform_test"
  parent_node_id = 2014254
  remark         = "for terraform test"
    tags = {
    "tf-test1" = "org_node_update"
  }
}

`
