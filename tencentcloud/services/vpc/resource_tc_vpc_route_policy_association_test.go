package vpc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudVpcRoutePolicyAssociationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRoutePolicyAssociation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_route_policy_association.example", "id"),
				),
			},
			{
				Config: testAccVpcRoutePolicyAssociationUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_route_policy_association.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_vpc_route_policy_association.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcRoutePolicyAssociation = `
resource "tencentcloud_vpc_route_policy_association" "example" {
  route_policy_id = "rrp-7dnu4yoi"
  route_table_id  = "rtb-389phpuq"
  priority        = 10
}
`

const testAccVpcRoutePolicyAssociationUpdate = `
resource "tencentcloud_vpc_route_policy_association" "example" {
  route_policy_id = "rrp-7dnu4yoi"
  route_table_id  = "rtb-389phpuq"
  priority        = 10
}
`
