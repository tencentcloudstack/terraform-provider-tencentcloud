provider "tencentcloud" {
  region = "ap-shanghai"
}

resource "tencentcloud_cdh_instance" "foo" {
  availability_zone = var.availability_zone
  host_type = "HM50"
  charge_type = "PREPAID"
  prepaid_period = 1
  host_name = "test"
  prepaid_renew_flag = "DISABLE_NOTIFY_AND_MANUAL_RENEW"
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

resource "tencentcloud_placement_group" "foo" {
  name = "test"
  type = "HOST"
}

resource "tencentcloud_instance" "foo" {
  availability_zone = var.availability_zone
  instance_name     = var.instance_name
  image_id          = var.image_id
  key_name          = tencentcloud_key_pair.random_key.id
  placement_group_id = tencentcloud_placement_group.foo.id
  security_groups               = [var.security_group_id]
  system_disk_type  = "CLOUD_PREMIUM"

  instance_charge_type = "CDHPAID"
  cdh_instance_type     = "CDH_10C10G"
  cdh_host_id = tencentcloud_cdh_instance.foo.id

  vpc_id                     = var.vpc_id
  subnet_id                  = var.subnet_id
  allocate_public_ip    = true
  internet_max_bandwidth_out = 2
  count                      = 3

  data_disks {
    data_disk_type = "CLOUD_PREMIUM"
    data_disk_size = 50
    encrypt = false
  }
}
