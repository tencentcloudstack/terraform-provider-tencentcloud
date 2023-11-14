package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudOrganizationOrgAuthNodeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgAuthNodeDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_organization_org_auth_node.org_auth_node")),
			},
		},
	})
}

const testAccOrganizationOrgAuthNodeDataSource = `

data "tencentcloud_organization_org_auth_node" "org_auth_node" {
  auth_name = &lt;nil&gt;
  }

`
