package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcNatDcRouteDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcNatDcRouteDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_nat_dc_route.nat_dc_route")),
			},
		},
	})
}

const testAccVpcNatDcRouteDataSource = `

data "tencentcloud_vpc_nat_dc_route" "nat_dc_route" {
  nat_gateway_id = "nat-b78el568"
  vpc_id = "vpc-34tb3e5j"
  }

`
