package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixEniIpv4AddressResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},

		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEniIpv4Address,
				Check: resource.ComposeTestCheckFunc(
					//testAccCheckBandwidthPackageAttachmentExists("tencentcloud_vpc_ipv6_eni_address"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni_ipv4_address.eni_ipv4_address", "network_interface_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni_ipv4_address.eni_ipv4_address", "private_ip_addresses.0.private_ip_address"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni_ipv4_address.eni_ipv4_address", "private_ip_addresses.0.primary"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni_ipv4_address.eni_ipv4_address", "private_ip_addresses.0.address_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni_ipv4_address.eni_ipv4_address", "private_ip_addresses.0.description"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni_ipv4_address.eni_ipv4_address", "private_ip_addresses.0.is_wan_ip_blocked"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni_ipv4_address.eni_ipv4_address", "private_ip_addresses.0.state"),
				),
			},
			{
				ResourceName:      "tencentcloud_eni_ipv4_address.eni_ipv4_address",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccEniIpv4Address = `

data "tencentcloud_enis" "eni" {
  name = "Primary ENI"
}

resource "tencentcloud_eni_ipv4_address" "eni_ipv4_address" {
  network_interface_id = data.tencentcloud_enis.eni.enis.0.id
  secondary_private_ip_address_count = 3
}

`
