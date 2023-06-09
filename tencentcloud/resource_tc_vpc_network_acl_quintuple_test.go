package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixVpcNetworkAclQuintupleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcNetworkAclQuintuple,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_network_acl_quintuple.network_acl_quintuple", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_network_acl_quintuple.network_acl_quintuple",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcNetworkAclQuintuple = `

resource "tencentcloud_vpc_network_acl_quintuple" "network_acl_quintuple" {
  network_acl_id = ""
  network_acl_quintuple_set {
		ingress {
			protocol = ""
			description = ""
			source_port = ""
			source_cidr = ""
			destination_port = ""
			destination_cidr = ""
			action = ""
			network_acl_quintuple_entry_id = ""
			priority = 
			create_time = ""
			network_acl_direction = ""
		}
		egress {
			protocol = ""
			description = ""
			source_port = ""
			source_cidr = ""
			destination_port = ""
			destination_cidr = ""
			action = ""
			network_acl_quintuple_entry_id = ""
			priority = 
			create_time = ""
			network_acl_direction = ""
		}

  }
}

`
