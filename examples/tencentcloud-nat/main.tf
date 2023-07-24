data "tencentcloud_availability_zones" "my_zones" {}

data "tencentcloud_images" "my_image" {
  os_name = "centos"
}

data "tencentcloud_instance_types" "my_instance_types" {
  cpu_core_count = 1
  memory_size    = 1
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_nat_gateway_vpc"
}

resource "tencentcloud_eip" "eip_example1" {
  name = "tf_nat_gateway_eip1"
}

resource "tencentcloud_eip" "eip_example2" {
  name = "tf_nat_gateway_eip2"
}

resource "tencentcloud_nat_gateway" "example" {
  name             = "tf_example_nat_gateway"
  vpc_id           = tencentcloud_vpc.vpc.id
  bandwidth        = 100
  max_concurrent   = 1000000
  assigned_eip_set = [
    tencentcloud_eip.eip_example1.public_ip,
    tencentcloud_eip.eip_example2.public_ip,
  ]
  tags = {
    tf_tag_key = "tf_tag_value"
  }
}

# Create route_table and entry
resource "tencentcloud_route_table" "my_route_table" {
  vpc_id = tencentcloud_vpc.vpc.id
  name   = "terraform test"
}

resource "tencentcloud_route_table_entry" "my_route_entry" {
  route_table_id         = tencentcloud_route_table.my_route_table.id
  destination_cidr_block = "10.0.0.0/8"
  next_type              = "NAT"
  next_hub               = tencentcloud_nat_gateway.example.id
}

# Create Subnet
resource "tencentcloud_subnet" "my_subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "terraform test"
  cidr_block        = "172.29.23.0/24"
  availability_zone = data.tencentcloud_availability_zones.my_zones.zones.0.name
  route_table_id    = tencentcloud_route_table.my_route_table.id
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
  vpc_id                     = tencentcloud_vpc.vpc.id
  subnet_id                  = tencentcloud_subnet.my_subnet.id
  internet_max_bandwidth_out = 20
}

# Add DNAT Entry
resource "tencentcloud_dnat" "dev_dnat" {
  vpc_id       = tencentcloud_nat_gateway.example.vpc_id
  nat_id       = tencentcloud_nat_gateway.example.id
  protocol     = "TCP"
  elastic_ip   = tencentcloud_eip.eip_dev_dnat.public_ip
  elastic_port = "80"
  private_ip   = tencentcloud_instance.my_instance.private_ip
  private_port = "9001"
}

resource "tencentcloud_dnat" "test_dnat" {
  vpc_id       = tencentcloud_vpc.vpc.id
  nat_id       = tencentcloud_nat_gateway.example.id
  protocol     = "UDP"
  elastic_ip   = tencentcloud_eip.eip_test_dnat.public_ip
  elastic_port = "8080"
  private_ip   = tencentcloud_instance.my_instance.private_ip
  private_port = "9002"
}

# Subnet Nat gateway snat
resource "tencentcloud_nat_gateway_snat" "my_subnet_snat" {
  nat_gateway_id    = tencentcloud_nat_gateway.example.id
  resource_type     = "SUBNET"
  subnet_id         = tencentcloud_subnet.my_subnet.id
  subnet_cidr_block = tencentcloud_subnet.my_subnet.cidr_block
  description       = "terraform test"
  public_ip_addr    = [
    tencentcloud_eip.eip_dev_dnat.public_ip,
    tencentcloud_eip.eip_test_dnat.public_ip,
  ]
}

# NetWorkInterface Nat gateway snat
resource "tencentcloud_nat_gateway_snat" "my_instance_snat" {
  nat_gateway_id           = tencentcloud_nat_gateway.example.id
  resource_type            = "NETWORKINTERFACE"
  instance_id              = tencentcloud_instance.my_instance.id
  instance_private_ip_addr = tencentcloud_instance.my_instance.private_ip
  description              = "terraform test"
  public_ip_addr           = [
    tencentcloud_eip.eip_dev_dnat.public_ip,
  ]
}
