package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcIpv6CidrBlockResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcIpv6CidrBlock,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_ipv6_cidr_block.ipv6_cidr_block", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_ipv6_cidr_block.ipv6_cidr_block",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcIpv6CidrBlock = `

resource "tencentcloud_vpc" "cidr-block" {
  name         = "ipv6-cidr-block-for-test"
  cidr_block   = "10.0.0.0/16"
  is_multicast = false
}

resource "tencentcloud_vpc_ipv6_cidr_block" "ipv6_cidr_block" {
  vpc_id = tencentcloud_vpc.cidr-block.id
}

`
