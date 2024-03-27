package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixEniIpv6AddressResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},

		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEniIpv6Address,
				Check: resource.ComposeTestCheckFunc(
					//testAccCheckBandwidthPackageAttachmentExists("tencentcloud_vpc_ipv6_eni_address"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni_ipv6_address.ipv6_eni_address", "network_interface_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni_ipv6_address.ipv6_eni_address", "ipv6_addresses.0.address"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni_ipv6_address.ipv6_eni_address", "ipv6_addresses.0.primary"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni_ipv6_address.ipv6_eni_address", "ipv6_addresses.0.address_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni_ipv6_address.ipv6_eni_address", "ipv6_addresses.0.description"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni_ipv6_address.ipv6_eni_address", "ipv6_addresses.0.is_wan_ip_blocked"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni_ipv6_address.ipv6_eni_address", "ipv6_addresses.0.state"),
				),
			},
			{
				ResourceName:      "tencentcloud_eni_ipv6_address.ipv6_eni_address",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccEniIpv6Address = `

data "tencentcloud_availability_zones" "zones" {}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example-ipv6"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
  name              = "subnet-example-ipv6"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_eni" "eni" {
  name        = "eni-example-ipv6"
  vpc_id      = tencentcloud_vpc.vpc.id
  subnet_id   = tencentcloud_subnet.subnet.id
  description = "eni desc."
  ipv4_count  = 1
}

resource "tencentcloud_vpc_ipv6_cidr_block" "example" {
  vpc_id = tencentcloud_vpc.vpc.id
}

resource "tencentcloud_vpc_ipv6_subnet_cidr_block" "example" {
  vpc_id = tencentcloud_vpc.vpc.id
  ipv6_subnet_cidr_blocks {
    subnet_id       = tencentcloud_subnet.subnet.id
    ipv6_cidr_block = "2402:4e00:1018:6700::/64"
  }
}

resource "tencentcloud_eni_ipv6_address" "ipv6_eni_address" {
  network_interface_id = tencentcloud_eni.eni.id
  ipv6_address_count   = 1
}

`
