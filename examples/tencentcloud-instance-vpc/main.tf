data "tencentcloud_images" "my_favorate_image" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos"
}

data "tencentcloud_instance_types" "my_favorate_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S1"]
  }

  cpu_core_count = 1
  memory_size    = 1
}

data "tencentcloud_availability_zones" "my_favorate_zones" {}

resource "tencentcloud_vpc" "my_vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_vpc_test"
}

resource "tencentcloud_subnet" "my_subnet" {
  vpc_id = tencentcloud_vpc.my_vpc.id

  //  vpc_id     = "vpc-csybef02"
  availability_zone = data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name
  name              = "tf_test_subnet"
  cidr_block        = "10.0.2.0/24"
}

resource "tencentcloud_instance" "instance-vpc-example" {
  instance_name     = var.instance_name
  availability_zone = data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name
  image_id          = data.tencentcloud_images.my_favorate_image.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"

  vpc_id    = tencentcloud_vpc.my_vpc.id
  subnet_id = tencentcloud_subnet.my_subnet.id

  internet_max_bandwidth_out = 1
}
