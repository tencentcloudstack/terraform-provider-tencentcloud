package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudOrganizationNodesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccOrganizationNodesDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_organization_nodes.organization_nodes"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_organization_nodes.organization_nodes", "items.#"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_organization_nodes.organization_nodes", "items.0.name"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_organization_nodes.organization_nodes", "items.0.node_id"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_organization_nodes.organization_nodes", "items.0.parent_node_id"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_organization_nodes.organization_nodes", "items.0.create_time"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_organization_nodes.organization_nodes", "items.0.update_time"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_organization_nodes.organization_nodes", "items.0.tags.#"),
			),
		}},
	})
}

const testAccOrganizationNodesDataSource = `

data "tencentcloud_organization_nodes" "organization_nodes" {
    tags {
        tag_key = "createBy"
        tag_value = "terraform"
    }
}
`
