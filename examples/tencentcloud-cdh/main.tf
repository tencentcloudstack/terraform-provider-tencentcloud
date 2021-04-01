data "tencentcloud_images" "my_favorate_image" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = var.os_name
}

resource "tencentcloud_cdh_instance" "foo" {
  availability_zone = var.availability_zone
  host_type = "HC20"
  charge_type = "PREPAID"
  prepaid_period = 1
  host_name = "test"
  prepaid_renew_flag = "NOTIFY_AND_MANUAL_RENEW"
}

data "tencentcloud_cdh_instances" "list" {
  availability_zone = var.availability_zone
  host_id = tencentcloud_cdh_instance.foo.id
  host_name = "test"
  host_state = "RUNNING"
}

resource "tencentcloud_key_pair" "random_key" {
  key_name   = "tf_example_key6"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDjd8fTnp7Dcuj4mLaQxf9Zs/ORgUL9fQxRCNKkPgP1paTy1I513maMX126i36Lxxl3+FUB52oVbo/FgwlIfX8hyCnv8MCxqnuSDozf1CD0/wRYHcTWAtgHQHBPCC2nJtod6cVC3kB18KeV4U7zsxmwFeBIxojMOOmcOBuh7+trRw=="
}

resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_zone
  image_id          = data.tencentcloud_images.my_favorate_image.images.0.image_id
  instance_charge_type = "CDHPAID"
  cdh_instance_type     = "CDH_10C10G"
  cdh_host_id = tencentcloud_cdh_instance.foo.id
  key_name          = tencentcloud_key_pair.random_key.id
  system_disk_type  = "CLOUD_PREMIUM"

  allocate_public_ip    = true
  internet_max_bandwidth_out = 2
  count                      = 1
}
