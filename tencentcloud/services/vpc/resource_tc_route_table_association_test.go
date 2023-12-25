package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudRouteTableAssociationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableAssociation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_route_table_association.route_table_association", "id"),
					resource.TestCheckResourceAttr("tencentcloud_route_table_association.route_table_association", "route_table_id", "rtb-5toos5sy"),
				),
			},
			{
				Config: testAccRouteTableAssociationUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_route_table_association.route_table_association", "id"),
					resource.TestCheckResourceAttr("tencentcloud_route_table_association.route_table_association", "route_table_id", "rtb-pp764dr4"),
				),
			},
			{
				ResourceName:      "tencentcloud_route_table_association.route_table_association",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRouteTableAssociation = `

resource "tencentcloud_route_table_association" "route_table_association" {
  route_table_id = "rtb-5toos5sy"
  subnet_id      = "subnet-2y2omd4k"
}

`

const testAccRouteTableAssociationUpdate = `

resource "tencentcloud_route_table_association" "route_table_association" {
  route_table_id = "rtb-pp764dr4"
  subnet_id      = "subnet-2y2omd4k"
}

`
