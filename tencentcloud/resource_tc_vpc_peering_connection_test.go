package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcPeeringConnectionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPeeringConnection,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_peering_connection.peering_connection", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_peering_connection.peering_connection",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcPeeringConnection = `

resource "tencentcloud_vpc_peering_connection" "peering_connection" {
  source_vpc_id = "vpc-abcdef"
  peering_connection_name = "name"
  destination_vpc_id = "vpc-abc1234"
  destination_uin = "12345678"
  destination_region = "ap-beijing"
  bandwidth = 100
  type = "VPC_PEER"
  charge_type = "POSTPAID_BY_DAY_MAX"
  qos_level = "AU"
}

`
