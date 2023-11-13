package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcIpv6AddressResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcIpv6Address,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_ipv6_address.ipv6_address", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_ipv6_address.ipv6_address",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcIpv6Address = `

resource "tencentcloud_vpc_ipv6_address" "ipv6_address" {
  ip6_addresses = 
  internet_max_bandwidth_out = 200
  internet_charge_type = "TRAFFIC_POSTPAID_BY_HOUR"
  bandwidth_package_id = "bwp-34rfgt56"
}

`
