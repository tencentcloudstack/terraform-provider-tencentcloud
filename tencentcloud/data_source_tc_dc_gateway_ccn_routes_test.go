package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTencentCloudDcgV3CcnRoutesInstancesBasic(t *testing.T) {
	t.Parallel()

	var rKey = "data.tencentcloud_dc_gateway_ccn_routes.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudDcgCcnRoutesInstances,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckTencentCloudDataSourceID(rKey),
					resource.TestCheckResourceAttrSet(rKey, "instance_list.#"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudDcgCcnRoutesInstances = `
resource "tencentcloud_ccn" "main" {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}

resource "tencentcloud_dc_gateway" "ccn_main" {
  name                = "ci-cdg-ccn-test"
  network_instance_id = tencentcloud_ccn.main.id
  network_type        = "CCN"
  gateway_type        = "NORMAL"
}

resource "tencentcloud_dc_gateway_ccn_route" "route1" {
  dcg_id     = tencentcloud_dc_gateway.ccn_main.id
  cidr_block = "10.1.1.0/32"
}

resource "tencentcloud_dc_gateway_ccn_route" "route2" {
  dcg_id     = tencentcloud_dc_gateway.ccn_main.id
  cidr_block = "192.1.1.0/32"
}

#You need to sleep for a few seconds because there is a cache on the server
data "tencentcloud_dc_gateway_ccn_routes"  "test" {
  dcg_id = tencentcloud_dc_gateway.ccn_main.id
}
`
