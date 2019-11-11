resource "tencentcloud_vpc" "vpc" {
  name       = "test-cfs-vpc"
  cidr_block = "10.2.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = "${tencentcloud_vpc.vpc.id}"
  name              = "test-cfs-subnet"
  cidr_block        = "10.2.11.0/24"
  availability_zone = "${var.availability_zone}"
}

resource "tencentcloud_cfs_access_group" "foo" {
  name = "test_cfs_access_rule"
}

resource "tencentcloud_cfs_access_rule" "foo" {
  access_group_id = "${tencentcloud_cfs_access_group.foo.id}"
  auth_client_ip  = "10.10.1.0/24"
  rw_permission   = "RW"
  user_permission = "root_squash"
  priority        = 1
}

resource "tencentcloud_cfs_file_system" "foo" {
  name              = "test_cfs_file_system"
  availability_zone = "${var.availability_zone}"
  access_group_id   = "${tencentcloud_cfs_access_group.foo.id}"
  protocol          = "${var.protocol}"
  vpc_id            = "${tencentcloud_vpc.vpc.id}"
  subnet_id         = "${tencentcloud_subnet.subnet.id}"
}

data "tencentcloud_cfs_file_systems" "file_systems" {
  file_system_id = "${tencentcloud_cfs_file_system.foo.id}"
  name           = "${tencentcloud_cfs_file_system.foo.name}"
}

data "tencentcloud_cfs_access_rules" "access_rules" {
  access_group_id = "${tencentcloud_cfs_access_group.foo.id}"
  access_rule_id  = "${tencentcloud_cfs_access_rule.foo.id}"
}

data "tencentcloud_cfs_access_groups" "access_groups" {
  access_group_id = "${tencentcloud_cfs_access_group.foo.id}"
  name            = "${tencentcloud_cfs_access_group.foo.name}"
}