package vpc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudVpcReplaceRoutesWithRoutePolicyConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcReplaceRoutesWithRoutePolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_replace_routes_with_route_policy_config.example", "id"),
				),
			},
			{
				Config: testAccVpcReplaceRoutesWithRoutePolicyConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_replace_routes_with_route_policy_config.example", "id"),
				),
			},
		},
	})
}

const testAccVpcReplaceRoutesWithRoutePolicyConfig = `
resource "tencentcloud_vpc_replace_routes_with_route_policy_config" "example" {
  route_table_id = "rtb-olsbhnyc"
  routes {
    route_item_id      = "rti-araogi5t"
    force_match_policy = true
  }
}
`

const testAccVpcReplaceRoutesWithRoutePolicyConfigUpdate = `
resource "tencentcloud_vpc_replace_routes_with_route_policy_config" "example" {
  route_table_id = "rtb-olsbhnyc"
  routes {
    route_item_id      = "rti-araogi5t"
    force_match_policy = false
  }
}
`
