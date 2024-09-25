package tco_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudOrganizationOrganizationResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrganization,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_instance.organization", "id")),
			},
			{
				ResourceName:      "tencentcloud_organization_instance.organization",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudOrganizationOrganizationResource_rootNodeName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrganizationWithRootNodeName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_organization_instance.organization", "id"),
					resource.TestCheckResourceAttr("tencentcloud_organization_instance.organization", "root_node_name", "root_node_name"),
				),
			},
			{
				Config: testAccOrganizationOrganizationWithRootNodeNameUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_organization_instance.organization", "id"),
					resource.TestCheckResourceAttr("tencentcloud_organization_instance.organization", "root_node_name", "root_node_name_update"),
				),
			},
			{
				ResourceName:      "tencentcloud_organization_instance.organization",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccOrganizationOrganization = `

resource "tencentcloud_organization_instance" "organization" {
}

`

const testAccOrganizationOrganizationWithRootNodeName = `
resource "tencentcloud_organization_instance" "organization" {
	root_node_name = "root_node_name"
}
`

const testAccOrganizationOrganizationWithRootNodeNameUpdate = `
resource "tencentcloud_organization_instance" "organization" {
	root_node_name = "root_node_name_update"
}
`
