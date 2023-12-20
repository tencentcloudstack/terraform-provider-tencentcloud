package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcPeerConnectAccecptOrRejectResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPeerConnectAccecptOrReject,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_peer_connect_accecpt_or_reject.peer_connect_accecpt_or_reject", "id")),
			},
		},
	})
}

const testAccVpcPeerConnectAccecptOrReject = `

resource "tencentcloud_vpc_peer_connect_accecpt_or_reject" "peer_connect_accecpt_or_reject" {
  peering_connection_id = "pcx-k74fvy2e"
}

`
