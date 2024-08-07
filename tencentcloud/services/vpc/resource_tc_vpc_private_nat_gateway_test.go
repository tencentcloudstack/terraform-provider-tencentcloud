package vpc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudVpcPrivateNatGatewayResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPrivateNatGateway,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_private_nat_gateway.private_nat_gateway", "id"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_private_nat_gateway.private_nat_gateway", "nat_gateway_name", "private-nat-gateway"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_private_nat_gateway.private_nat_gateway", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_private_nat_gateway.private_nat_gateway", "cross_domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_private_nat_gateway.private_nat_gateway", "vpc_type"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_private_nat_gateway.private_nat_gateway", "ccn_id", ""),
				),
			},
			{
				Config: testAccVpcPrivateNatGatewayUpdateName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_private_nat_gateway.private_nat_gateway", "id"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_private_nat_gateway.private_nat_gateway", "nat_gateway_name", "private-nat-gateway-update"),
				),
			},
			{
				ResourceName:      "tencentcloud_vpc_private_nat_gateway.private_nat_gateway",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcPrivateNatGateway = `
resource "tencentcloud_vpc" "foo" {
  name       = "private-nat-gateway-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_vpc_private_nat_gateway" "private_nat_gateway" {
  nat_gateway_name = "private-nat-gateway"
  vpc_id = tencentcloud_vpc.foo.id
}
`

const testAccVpcPrivateNatGatewayUpdateName = `
resource "tencentcloud_vpc" "foo" {
  name       = "private-nat-gateway-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_vpc_private_nat_gateway" "private_nat_gateway" {
  nat_gateway_name = "private-nat-gateway-update"
  vpc_id = tencentcloud_vpc.foo.id
}
`
