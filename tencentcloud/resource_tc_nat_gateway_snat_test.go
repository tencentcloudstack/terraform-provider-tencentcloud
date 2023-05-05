package tencentcloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
					testAccCheckNatGatewaySnatExists("tencentcloud_nat_gateway_snat.my_subnet_snat"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway_snat.my_subnet_snat", "resource_type", "SUBNET"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway_snat.my_subnet_snat", "description", "terraform test"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway_snat.my_subnet_snat", "public_ip_addr.#", "2"),
				),
			},
			{
				Config: testAccNatGatewaySnatConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewaySnatExists("tencentcloud_nat_gateway_snat.my_subnet_snat"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway_snat.my_subnet_snat", "resource_type", "SUBNET"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway_snat.my_subnet_snat", "description", "terraform test2"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway_snat.my_subnet_snat", "public_ip_addr.#", "1"),
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

		if rs.Primary.ID == "" {
			return fmt.Errorf("nat gateway snat id is not set")
		}

		ids := strings.Split(rs.Primary.ID, FILED_SP)

		err, result := service.DescribeNatGatewaySnats(contextNil, ids[0], nil)
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

		ids := strings.Split(rs.Primary.ID, FILED_SP)

		service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		err, result := service.DescribeNatGatewaySnats(contextNil, ids[0], nil)
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

data "tencentcloud_vpc_instances" "my_vpc" {
  name = "Default-VPC"
}

data "tencentcloud_vpc_subnets" "my_subnet" {
  vpc_id    = data.tencentcloud_vpc_instances.my_vpc.instance_list.0.vpc_id
  subnet_id = "subnet-4o0zd840"
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
  vpc_id         = data.tencentcloud_vpc_instances.my_vpc.instance_list.0.vpc_id
  name           = "terraform test"
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
`
const testAccNatGatewaySnatConfigUpdate = `

data "tencentcloud_vpc_instances" "my_vpc" {
  name = "Default-VPC"
}

data "tencentcloud_vpc_subnets" "my_subnet" {
  vpc_id    = data.tencentcloud_vpc_instances.my_vpc.instance_list.0.vpc_id
  subnet_id = "subnet-4o0zd840"
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
  vpc_id         = data.tencentcloud_vpc_instances.my_vpc.instance_list.0.vpc_id
  name           = "terraform test"
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
  description       = "terraform test2"
  public_ip_addr = [
    tencentcloud_eip.eip_dev_dnat.public_ip,
  ]
}
`
