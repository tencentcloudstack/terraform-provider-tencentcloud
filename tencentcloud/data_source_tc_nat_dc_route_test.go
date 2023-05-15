package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixNatDcRouteDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNatDcRouteDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_nat_dc_route.nat_dc_route")),
			},
		},
	})
}

const testAccNatDcRouteDataSource = `

data "tencentcloud_nat_dc_route" "nat_dc_route" {
  nat_gateway_id = "nat-gnxkey2e"
  vpc_id         = "vpc-pyyv5k3v"
}
`
