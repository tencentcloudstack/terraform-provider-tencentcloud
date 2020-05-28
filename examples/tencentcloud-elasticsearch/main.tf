resource "tencentcloud_vpc" "vpc" {
  name       = "tf-es-vpc"
  cidr_block = "10.2.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "tf-es-subnet"
  cidr_block        = "10.2.11.0/24"
  availability_zone = var.availability_zone
}

resource "tencentcloud_elasticsearch_instance" "foo" {
  instance_name       = var.instance_name
  availability_zone   = var.availability_zone
  version             = "7.5.1"
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = tencentcloud_subnet.subnet.id
  password            = "Test12345"
  license_type        = "basic"
  basic_security_type = 2

  node_info_list {
    node_num  = 2
    node_type = var.instance_type
  }

  tags = {
    test = "test"
  }
}