data "tencentcloud_image" "my_favorate_image" {
  os_name = "centos"

  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

data "tencentcloud_instance_types" "my_favorate_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S2"]
  }

  cpu_core_count = 2
  memory_size    = 4
}

data "tencentcloud_availability_zones" "my_favorate_zones" {}

resource "tencentcloud_instance" "instance-without-specified-image-id-example" {
  instance_name = "${var.instance_name}"
  availability_zone     = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
  image_id      = "${data.tencentcloud_image.my_favorate_image.image_id}"
  instance_type = "${data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type}"
}

resource "tencentcloud_cbs_storage" "my-storage" {
  storage_type = "cloudBasic"
  storage_size = 10
  period       = 1
  availability_zone    = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
  storage_name = "${var.storage_name}"
}

resource "tencentcloud_cbs_storage_attachment" "my-attachment" {
  storage_id  = "${tencentcloud_cbs_storage.my-storage.id}"
  instance_id = "${tencentcloud_instance.instance-without-specified-image-id-example.id}"
}
