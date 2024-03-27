package vpc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcEniIpv4AddressResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEniIpv4Address,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_eni_ipv4_address.eni_ipv4_address", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_eni_ipv4_address.eni_ipv4_address",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcEniIpv4Address = `

resource "tencentcloud_vpc_eni_ipv4_address" "eni_ipv4_address" {
  network_interface_id = ""
  private_ip_addresses {
		private_ip_address = ""
		primary = 
		public_ip_address = ""
		address_id = ""
		description = ""
		is_wan_ip_blocked = 
		state = ""
		qos_level = ""

  }
  secondary_private_ip_address_count = 
  qos_level = ""
}

`
