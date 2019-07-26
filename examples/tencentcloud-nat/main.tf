data "tencentcloud_availability_zones" "my_favorate_zones" {}

data "tencentcloud_image" "my_favorate_image" {
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

# Create VPC and Subnet
resource "tencentcloud_vpc" "main" {
  name       = "terraform test"
  cidr_block = "10.6.0.0/16"
}

resource "tencentcloud_subnet" "main_subnet" {
  vpc_id            = "${tencentcloud_vpc.main.id}"
  name              = "terraform test subnet"
  cidr_block        = "10.6.7.0/24"
  availability_zone = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
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
  vpc_id         = "${tencentcloud_vpc.main.id}"
  name           = "terraform test"
  max_concurrent = 3000000
  bandwidth      = 500

  assigned_eip_set = [
    "${tencentcloud_eip.eip_dev_dnat.public_ip}",
    "${tencentcloud_eip.eip_test_dnat.public_ip}",
  ]
}

# Create CVM
resource "tencentcloud_instance" "foo" {
  availability_zone = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
  image_id          = "${data.tencentcloud_image.my_favorate_image.image_id}"
  vpc_id            = "${tencentcloud_vpc.main.id}"
  subnet_id         = "${tencentcloud_subnet.main_subnet.id}"
}

# Add DNAT Entry
resource "tencentcloud_dnat" "dev_dnat" {
  vpc_id       = "${tencentcloud_nat_gateway.my_nat.vpc_id}"
  nat_id       = "${tencentcloud_nat_gateway.my_nat.id}"
  protocol     = "tcp"
  elastic_ip   = "${tencentcloud_eip.eip_dev_dnat.public_ip}"
  elastic_port = "80"
  private_ip   = "${tencentcloud_instance.foo.private_ip}"
  private_port = "9001"
}

resource "tencentcloud_dnat" "test_dnat" {
  vpc_id       = "${tencentcloud_nat_gateway.my_nat.vpc_id}"
  nat_id       = "${tencentcloud_nat_gateway.my_nat.id}"
  protocol     = "udp"
  elastic_ip   = "${tencentcloud_eip.eip_test_dnat.public_ip}"
  elastic_port = "8080"
  private_ip   = "${tencentcloud_instance.foo.private_ip}"
  private_port = "9002"
}
