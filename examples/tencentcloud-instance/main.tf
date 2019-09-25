data "tencentcloud_image" "my_favorate_image" {
  os_name = "${var.os_name}"

  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
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

resource "tencentcloud_key_pair" "random_key" {
  key_name   = "tf_example_key6"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDjd8fTnp7Dcuj4mLaQxf9Zs/ORgUL9fQxRCNKkPgP1paTy1I513maMX126i36Lxxl3+FUB52oVbo/FgwlIfX8hyCnv8MCxqnuSDozf1CD0/wRYHcTWAtgHQHBPCC2nJtod6cVC3kB18KeV4U7zsxmwFeBIxojMOOmcOBuh7+trRw=="
}

resource "tencentcloud_instance" "instance-without-specified-image-id-example" {
  instance_name     = "${var.instance_name}"
  availability_zone = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
  image_id          = "${data.tencentcloud_image.my_favorate_image.image_id}"
  instance_type     = "${data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type}"
  key_name          = "${tencentcloud_key_pair.random_key.id}"

  //  instance_charge_type                = "PREPAID"
  //  instance_charge_type_prepaid_period = 1

  disable_monitor_service    = true
  internet_max_bandwidth_out = 2
  count                      = 1
}
