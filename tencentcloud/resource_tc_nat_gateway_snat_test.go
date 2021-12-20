package tencentcloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudNatGatewaySnat_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNatGatewaySnatDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatGatewaySnatConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewaySnatExists("tencentcloud_nat_gateway.my_subnet_snat"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_subnet_snat", "resource_type", "SUBNET"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_subnet_snat", "description", "terraform test"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_subnet_snat", "public_ip_addr.#", "2"),
				),
			},
			{
				Config: testAccNatGatewaySnatConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewaySnatExists("tencentcloud_nat_gateway.my_subnet_snat"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_subnet_snat", "resource_type", "SUBNET"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_subnet_snat", "description", "terraform test2"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_subnet_snat", "public_ip_addr.#", "1"),
				),
			},
		},
	})
}

func testAccCheckNatGatewaySnatDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)

	service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_nat_gateway_snat" {
			continue
		}
		err, result := service.DescribeNatGatewaySnats(contextNil, rs.Primary.ID, nil)
		if err != nil {
			log.Printf("[CRITAL]%s read nat gateway snat failed, reason:%s\n ", logId, err.Error())
			return err
		}
		if len(result) != 0 {
			return fmt.Errorf("nat gateway snat id is still exists")
		}

	}
	return nil
}

func testAccCheckNatGatewaySnatExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("nat gateway snat instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("nat gateway snat id is not set")
		}
		service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		err, result := service.DescribeNatGatewaySnats(contextNil, rs.Primary.ID, nil)
		if err != nil {
			log.Printf("[CRITAL]%s read nat gateway snat failed, reason:%s\n ", logId, err.Error())
			return err
		}
		if len(result) != 1 {
			return fmt.Errorf("nat gateway snat id is not found")
		}
		return nil
	}
}

const testAccNatGatewaySnatConfig = `
data "tencentcloud_availability_zones" "my_zones" {}

data "tencentcloud_vpc" "my_vpc" {
  name = "Default-VPC"
}

data "tencentcloud_images" "my_image" {
  os_name = "centos"
}

data "tencentcloud_instance_types" "my_instance_types" {
  cpu_core_count = 1
  memory_size    = 1
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
  vpc_id         = data.tencentcloud_vpc.my_vpc.id
  name           = "terraform test"
  max_concurrent = 3000000
  bandwidth      = 500

  assigned_eip_set = [
    tencentcloud_eip.eip_dev_dnat.public_ip,
    tencentcloud_eip.eip_test_dnat.public_ip,
  ]
}

# Create route_table and entry
resource "tencentcloud_route_table" "my_route_table" {
  vpc_id = data.tencentcloud_vpc.my_vpc.id
  name   = "terraform test"
}
resource "tencentcloud_route_table_entry" "my_route_entry" {
  route_table_id         = tencentcloud_route_table.my_route_table.id
  destination_cidr_block = "10.0.0.0/8"
  next_type              = "NAT"
  next_hub               = tencentcloud_nat_gateway.my_nat.id
}

# Create Subnet
resource "tencentcloud_subnet" "my_subnet" {
  vpc_id            = data.tencentcloud_vpc.my_vpc.id
  name              = "terraform test"
  cidr_block        = "172.29.23.0/24"
  availability_zone = data.tencentcloud_availability_zones.my_zones.zones.0.name
  route_table_id    = tencentcloud_route_table.my_route_table.id
}

# Subnet Nat gateway snat
resource "tencentcloud_nat_gateway_snat" "my_subnet_snat" {
  nat_gateway_id    = tencentcloud_nat_gateway.my_nat.id
  resource_type     = "SUBNET"
  subnet_id         = tencentcloud_subnet.my_subnet.id
  subnet_cidr_block = tencentcloud_subnet.my_subnet.cidr_block
  description       = "terraform test"
  public_ip_addr = [
    tencentcloud_eip.eip_dev_dnat.public_ip,
    tencentcloud_eip.eip_test_dnat.public_ip,
  ]
}
`
const testAccNatGatewaySnatConfigUpdate = `
data "tencentcloud_availability_zones" "my_zones" {}

data "tencentcloud_vpc" "my_vpc" {
  name = "Default-VPC"
}

data "tencentcloud_images" "my_image" {
  os_name = "centos"
}

data "tencentcloud_instance_types" "my_instance_types" {
  cpu_core_count = 1
  memory_size    = 1
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
  vpc_id         = data.tencentcloud_vpc.my_vpc.id
  name           = "terraform test"
  max_concurrent = 3000000
  bandwidth      = 500

  assigned_eip_set = [
    tencentcloud_eip.eip_dev_dnat.public_ip,
    tencentcloud_eip.eip_test_dnat.public_ip,
  ]
}

# Create route_table and entry
resource "tencentcloud_route_table" "my_route_table" {
  vpc_id = data.tencentcloud_vpc.my_vpc.id
  name   = "terraform test"
}
resource "tencentcloud_route_table_entry" "my_route_entry" {
  route_table_id         = tencentcloud_route_table.my_route_table.id
  destination_cidr_block = "10.0.0.0/8"
  next_type              = "NAT"
  next_hub               = tencentcloud_nat_gateway.my_nat.id
}

# Create Subnet
resource "tencentcloud_subnet" "my_subnet" {
  vpc_id            = data.tencentcloud_vpc.my_vpc.id
  name              = "terraform test"
  cidr_block        = "172.29.23.0/24"
  availability_zone = data.tencentcloud_availability_zones.my_zones.zones.0.name
  route_table_id    = tencentcloud_route_table.my_route_table.id
}

# Subnet Nat gateway snat
resource "tencentcloud_nat_gateway_snat" "my_subnet_snat" {
  nat_gateway_id    = tencentcloud_nat_gateway.my_nat.id
  resource_type     = "SUBNET"
  subnet_id         = tencentcloud_subnet.my_subnet.id
  subnet_cidr_block = tencentcloud_subnet.my_subnet.cidr_block
  description       = "terraform test2"
  public_ip_addr = [
    tencentcloud_eip.eip_dev_dnat.public_ip,
  ]
}
`
