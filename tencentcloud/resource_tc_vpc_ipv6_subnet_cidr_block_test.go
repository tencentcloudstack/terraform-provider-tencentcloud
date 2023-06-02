package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixVpcIpv6SubnetCidrBlockResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcIpv6SubnetCidrBlock,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_ipv6_subnet_cidr_block.ipv6_subnet_cidr_block", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_ipv6_subnet_cidr_block.ipv6_subnet_cidr_block",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcIpv6SubnetCidrBlock = `

resource "tencentcloud_vpc_ipv6_subnet_cidr_block" "ipv6_subnet_cidr_block" {
  vpc_id = "vpc-7w3kgnpl"
  ipv6_subnet_cidr_blocks {
    subnet_id = "subnet-plg028y8"
    ipv6_cidr_block = "2402:4e00:1019:6a7b::/64"
  }
}

`
