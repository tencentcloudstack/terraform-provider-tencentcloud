package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcDhcpIpResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcDhcpIp,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_dhcp_ip.dhcp_ip", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_dhcp_ip.dhcp_ip",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcDhcpIp = `

resource "tencentcloud_vpc_dhcp_ip" "dhcp_ip" {
  vpc_id       = "vpc-86v957zb"
  subnet_id    = "subnet-enm92y0m"
  dhcp_ip_name = "terraform-test"
}

`
