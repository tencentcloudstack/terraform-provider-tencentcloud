package tco_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixOrganizationServiceAssign_basic -v
func TestAccTencentCloudNeedFixOrganizationServiceAssign_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_SMS) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationServiceAssign,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_organization_service_assign.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_organization_service_assign.example", "service_id", "15"),
					resource.TestCheckResourceAttr("tencentcloud_organization_service_assign.example", "management_scope", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_service_assign.example", "member_uins"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_service_assign.example", "management_scope_uins"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_service_assign.example", "management_scope_node_ids"),
				),
			},
			{
				ResourceName:      "tencentcloud_organization_service_assign.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccOrganizationServiceAssign = `
resource "tencentcloud_organization_service_assign" "example" {
  service_id       = 15
  management_scope = 2
  member_uins = [100013415241, 100078908111]
  management_scope_uins = [100019287759, 100020537485]
  management_scope_node_ids = [2024256, 2024259]
}
`
