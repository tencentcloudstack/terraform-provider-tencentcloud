package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixNatRefreshNatDcRouteResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNatRefreshNatDcRoute,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_nat_refresh_nat_dc_route.refresh_nat_dc_route", "id")),
			},
			{
				ResourceName:      "tencentcloud_nat_refresh_nat_dc_route.refresh_nat_dc_route",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccNatRefreshNatDcRoute = `

resource "tencentcloud_nat_refresh_nat_dc_route" "refresh_nat_dc_route" {
  nat_gateway_id = "nat-gnxkey2e"
  vpc_id         = "vpc-pyyv5k3v"
  dry_run = true
}

`
