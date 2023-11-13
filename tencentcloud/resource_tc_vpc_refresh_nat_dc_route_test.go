package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcRefreshNatDcRouteResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRefreshNatDcRoute,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_refresh_nat_dc_route.refresh_nat_dc_route", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_refresh_nat_dc_route.refresh_nat_dc_route",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcRefreshNatDcRoute = `

resource "tencentcloud_vpc_refresh_nat_dc_route" "refresh_nat_dc_route" {
  vpc_id = "vpc-34tb3e5j"
  nat_gateway_id = "nat-b78el568"
  dry_run = True
}

`
