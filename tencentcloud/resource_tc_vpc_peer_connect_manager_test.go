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

data "tencentcloud_user_info" "info" {}

locals {
  owner_uin = data.tencentcloud_user_info.info.owner_uin
}

resource "tencentcloud_vpc" "vpc" {
  name       = "tf-example-pcx"
  cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_vpc" "des_vpc" {
  name       = "tf-example-pcx-des"
  cidr_block = "172.16.0.0/16"
}
resource "tencentcloud_vpc_peer_connect_manager" "peer_connect_manager" {
  source_vpc_id = tencentcloud_vpc.vpc.id
  peering_connection_name = "example-iac"
  destination_vpc_id = tencentcloud_vpc.des_vpc.id
  destination_uin = local.owner_uin
  destination_region = "ap-guangzhou"
}


`
