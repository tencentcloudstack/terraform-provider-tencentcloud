data "tencentcloud_instance_types" "my_favorate_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S1"]
  }

  cpu_core_count = 1
  memory_size    = 1
}

data "tencentcloud_availability_zones" "my_favorate_zones" {}

resource "tencentcloud_instance" "jilei" {
  availability_zone = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
  instance_type     = "${data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type}"
  instance_name     = "tf_example_os_reinstall"

  disable_monitor_service = true

  image_id = "img-871lthrb" //  image_id = "img-31tjrtph"
  password = "test12345"    //  password = "test1234"
}
