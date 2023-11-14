package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudOrganizationQuitOrganizationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationQuitOrganization,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_quit_organization.quit_organization", "id")),
			},
			{
				ResourceName:      "tencentcloud_organization_quit_organization.quit_organization",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccOrganizationQuitOrganization = `

resource "tencentcloud_organization_quit_organization" "quit_organization" {
  org_id = &lt;nil&gt;
}

`
