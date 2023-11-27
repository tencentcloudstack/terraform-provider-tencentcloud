package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudRouteTableAssociationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableAssociation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_route_table_association.route_table_association", "id")),
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
