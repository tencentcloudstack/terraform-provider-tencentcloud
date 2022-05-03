package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNatGatewaySnatsDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudNatGatewaySnatsDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_nat_gateway_snats.snat"),
					resource.TestCheckResourceAttr("data.tencentcloud_nat_gateway_snats.snat", "nats.#", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_nat_gateway_snats.snat", "nats.0.resource_type", "SUBNET"),
					resource.TestCheckResourceAttr("data.tencentcloud_nat_gateway_snats.snat", "nats.1.resource_type", "NETWORKINTERFACE"),
					resource.TestCheckResourceAttr("data.tencentcloud_nat_gateway_snats.snat", "nats.0.description", "terraform test"),
					resource.TestCheckResourceAttr("data.tencentcloud_nat_gateway_snats.snat", "nats.1.description", "terraform test"),
				),
			},
		},
	})
}

const testAccTencentCloudNatGatewaySnatsDataSourceConfig_basic = `
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
  availability_zone = data.tencentcloud_availability_zones.my_zones.zones.0.name 
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
  cidr_block        = "172.16.128.0/20"
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

# Create instance
resource "tencentcloud_instance" "my_instance" {
  instance_name              = "terraform test"
  availability_zone          = data.tencentcloud_availability_zones.my_zones.zones.0.name
  image_id                   = data.tencentcloud_images.my_image.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.my_instance_types.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  system_disk_size           = 50
  hostname                   = "user"
  project_id                 = 0
  vpc_id                     = data.tencentcloud_vpc.my_vpc.id
  subnet_id                  = tencentcloud_subnet.my_subnet.id
  //internet_max_bandwidth_out = 20
}

# NetWorkInterface Nat gateway snat
resource "tencentcloud_nat_gateway_snat" "my_instance_snat" {
  nat_gateway_id           = tencentcloud_nat_gateway.my_nat.id
  resource_type            = "NETWORKINTERFACE"
  instance_id              = tencentcloud_instance.my_instance.id
  instance_private_ip_addr = tencentcloud_instance.my_instance.private_ip
  description              = "terraform test"
  public_ip_addr = [
    tencentcloud_eip.eip_dev_dnat.public_ip,
  ]
}

data "tencentcloud_nat_gateway_snats" "snat" {
  nat_gateway_id     = tencentcloud_nat_gateway.my_nat.id
}
`
