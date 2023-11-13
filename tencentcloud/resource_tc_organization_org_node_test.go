package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudOrganizationOrgNodeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgNode,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_org_node.org_node", "id")),
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
  node_id = &lt;nil&gt;
  parent_node_id = &lt;nil&gt;
  name = &lt;nil&gt;
  remark = &lt;nil&gt;
  create_time = &lt;nil&gt;
  update_time = &lt;nil&gt;
}

`
