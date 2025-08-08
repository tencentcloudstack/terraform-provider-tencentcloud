package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudOrganizationResourceToShareMemberDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccOrganizationResourceToShareMemberDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_organization_resource_to_share_member.example"),
				resource.TestCheckResourceAttr("data.tencentcloud_organization_resource_to_share_member.example", "items.#", "1"),
			),
		}},
	})
}

const testAccOrganizationResourceToShareMemberDataSource = `
data "tencentcloud_organization_resource_to_share_member" "example" {
  area                 = ""
  search_key           = ""
  type                 = ""
  product_resource_ids = []
}
`
