package vpc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudVpcRoutePolicyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRoutePolicy,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_route_policy.example", "id"),
				),
			},
			{
				Config: testAccVpcRoutePolicyUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_route_policy.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_vpc_route_policy.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcRoutePolicy = `
resource "tencentcloud_vpc_route_policy" "example" {
  route_policy_name        = "tf-example"
  route_policy_description = "remark."
}
`

const testAccVpcRoutePolicyUpdate = `
resource "tencentcloud_vpc_route_policy" "example" {
  route_policy_name        = "tf-example-update"
  route_policy_description = "remark."
}
`
