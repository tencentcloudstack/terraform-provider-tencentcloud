package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudOrganizationQuitOrganizationOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationQuitOrganizationOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_organization_quit_organization_operation.quit_organization_operation", "id"),
					resource.TestCheckResourceAttr("tencentcloud_organization_quit_organization_operation.quit_organization_operation", "org_id", "45155")),
			},
			{
				ResourceName:      "tencentcloud_organization_quit_organization_operation.quit_organization_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccOrganizationQuitOrganizationOperation = `

resource "tencentcloud_organization_quit_organization_operation" "quit_organization_operation" {
  org_id = 45155
}

`
