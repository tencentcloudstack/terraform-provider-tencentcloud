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

  cpu_core_count = 1
  memory_size    = 2
}

data "tencentcloud_availability_zones" "my_favorate_zones" {}

resource "tencentcloud_instance" "my_instance" {
  instance_name     = "terraform_automation_test_kuruk"
  availability_zone = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
  image_id          = "${data.tencentcloud_image.my_favorate_image.image_id}"
  instance_type     = "${data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type}"

  system_disk_type = "CLOUD_BASIC"

  data_disks = [
    {
      data_disk_type = "CLOUD_BASIC"
      data_disk_size = 70
    },
  ]

  disable_security_service = true
  disable_monitor_service  = true
}

resource "tencentcloud_eip" "my_eip" {
  name = "tf_auto_test"
}

resource "tencentcloud_eip_association" "foo" {
  eip_id      = "${tencentcloud_eip.my_eip.id}"
  instance_id = "${tencentcloud_instance.my_instance.id}"
}
