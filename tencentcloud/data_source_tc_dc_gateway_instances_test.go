package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTencentCloudDcgV3InstancesBasic(t *testing.T) {

	var nameKey = "data.tencentcloud_dc_gateway_instances.name_select"
	var idKey = "data.tencentcloud_dc_gateway_instances.id_select"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudDcgInstances,
				Check: resource.ComposeTestCheckFunc(

					//name filter
					testAccCheckTencentCloudDataSourceID(nameKey),
					resource.TestCheckResourceAttrSet(nameKey, "instance_list.#"),

					//id filter
					testAccCheckTencentCloudDataSourceID(idKey),
					resource.TestCheckResourceAttr(idKey, "instance_list.#", "1"),

					resource.TestCheckResourceAttr(idKey, "instance_list.0.gateway_type", "NORMAL"),
					resource.TestCheckResourceAttr(idKey, "instance_list.0.name", "ci-cdg-ccn-test"),
					resource.TestCheckResourceAttr(idKey, "instance_list.0.network_type", "CCN"),

					resource.TestCheckResourceAttrSet(idKey, "instance_list.0.network_instance_id"),
					resource.TestCheckResourceAttrSet(idKey, "instance_list.0.create_time"),
					resource.TestCheckResourceAttrSet(idKey, "instance_list.0.cnn_route_type"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudDcgInstances = `
resource "tencentcloud_ccn" "main" {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}

resource "tencentcloud_dc_gateway" "ccn_main" {
  name                = "ci-cdg-ccn-test"
  network_instance_id = "${tencentcloud_ccn.main.id}"
  network_type        = "CCN"
  gateway_type        = "NORMAL"
}

data "tencentcloud_dc_gateway_instances" "name_select"{
  name = "${tencentcloud_dc_gateway.ccn_main.name}"
}

data "tencentcloud_dc_gateway_instances"  "id_select" {
  dcg_id = "${tencentcloud_dc_gateway.ccn_main.id}"
}
`
