package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudCcnRoutesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcCcnRoutes,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ccn_routes.ccn_routes", "id")),
			},
			{
				ResourceName:      "tencentcloud_ccn_routes.ccn_routes",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcCcnRoutes = `

resource "tencentcloud_ccn_routes" "ccn_routes" {
  ccn_id = "ccn-39lqkygf"
  route_id = "ccnr-3o0dfyuw"
  switch = "off"
}

`
