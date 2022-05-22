package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixNatGatewaySnatsDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudNatGatewaySnatsDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_nat_gateway_snats.snat"),
					resource.TestCheckResourceAttr("data.tencentcloud_nat_gateway_snats.snat", "snat_list.#", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_nat_gateway_snats.snat", "snat_list.0.resource_type", "SUBNET"),
					resource.TestCheckResourceAttr("data.tencentcloud_nat_gateway_snats.snat", "snat_list.1.resource_type", "NETWORKINTERFACE"),
					resource.TestCheckResourceAttr("data.tencentcloud_nat_gateway_snats.snat", "snat_list.0.description", "terraform test"),
					resource.TestCheckResourceAttr("data.tencentcloud_nat_gateway_snats.snat", "snat_list.1.description", "terraform test"),
				),
			},
		},
	})
}

const testAccTencentCloudNatGatewaySnatsDataSourceConfig_basic = `
data "tencentcloud_vpc_instances" "my_vpc" {
  name = "Default-VPC"
}

data "tencentcloud_vpc_subnets" "my_subnet" {
  vpc_id    = data.tencentcloud_vpc_instances.my_vpc.instance_list.0.vpc_id
  subnet_id = "subnet-4o0zd840"
}

# Create EIP
resource "tencentcloud_eip" "eip_dev_dnat" {
  name = "terraform_test"
}
resource "tencentcloud_eip" "eip_test_dnat" {
  name = "terraform_test"
}

# Create NAT Gateway
resource "tencentcloud_nat_gateway" "my_nat" {
  vpc_id         = data.tencentcloud_vpc_instances.my_vpc.instance_list.0.vpc_id
  name           = "terraform datasource test"
  max_concurrent = 3000000
  bandwidth      = 500

  assigned_eip_set = [
    tencentcloud_eip.eip_dev_dnat.public_ip,
    tencentcloud_eip.eip_test_dnat.public_ip,
  ]
}

# Subnet Nat gateway snat
resource "tencentcloud_nat_gateway_snat" "my_subnet_snat" {
  nat_gateway_id    = tencentcloud_nat_gateway.my_nat.id
  resource_type     = "SUBNET"
  subnet_id         = data.tencentcloud_vpc_subnets.my_subnet.instance_list.0.subnet_id
  subnet_cidr_block = data.tencentcloud_vpc_subnets.my_subnet.instance_list.0.cidr_block
  description       = "terraform test"
  public_ip_addr = [
    tencentcloud_eip.eip_dev_dnat.public_ip,
    tencentcloud_eip.eip_test_dnat.public_ip,
  ]
}

data "tencentcloud_nat_gateway_snats" "snat" {
  nat_gateway_id     = tencentcloud_nat_gateway.my_nat.id
}
`
