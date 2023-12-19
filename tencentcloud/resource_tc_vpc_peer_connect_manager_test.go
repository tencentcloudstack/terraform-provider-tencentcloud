package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcPeerConnectManagerResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPeerConnectManager,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_peer_connect_manager.peer_connect_manager", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_peer_connect_manager.peer_connect_manager",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcPeerConnectManager = `

resource "tencentcloud_vpc_peer_connect_manager" "peer_connect_manager" {
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
