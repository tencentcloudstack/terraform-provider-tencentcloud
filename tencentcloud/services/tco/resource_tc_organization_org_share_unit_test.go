package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudOrganizationOrgShareUnitResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgShareUnit,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_org_share_unit.org_share_unit", "id")),
			},
			{
				ResourceName:      "tencentcloud_organization_org_share_unit.org_share_unit",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccOrganizationOrgShareUnit = `

resource "tencentcloud_organization_org_share_unit" "org_share_unit" {
  name = "iac-test"
  area = "ap-guangzhou"
  description = "iac-test"
}

`
