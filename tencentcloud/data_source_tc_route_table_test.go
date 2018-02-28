package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceTencentCloudRouteTable_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTencentCloudRouteTableConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_route_table.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_route_table.foo", "name", "ci-terraform-routetable-do-not-delete"),
				),
			},
		},
	})
}

const testAccDataSourceTencentCloudRouteTableConfig = `
data "tencentcloud_route_table" "foo" {
	route_table_id = "rtb-o0beqbrc"
}
`
