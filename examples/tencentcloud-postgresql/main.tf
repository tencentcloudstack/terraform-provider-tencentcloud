resource "tencentcloud_vpc" "foo" {
  name       = "example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "foo" {
  name              = "example"
  availability_zone = var.availability_zone
  vpc_id            = tencentcloud_vpc.foo.id
  cidr_block        = "10.0.0.0/24"
  is_multicast      = false
}

resource "tencentcloud_postgresql_instance" "example" {
  name              = "tf_postsql_instance_111"
  availability_zone = var.availability_zone
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.foo.id
  subnet_id         = tencentcloud_subnet.foo.id
  engine_version    = "9.3.5"
  root_user         = "root"
  root_password     = "1qaA2k1wgvfa3ZZZ"
  charset           = "UTF8"
  project_id        = 0
  memory            = 2
  storage           = 10
}

data "tencentcloud_postgresql_instances" "id_example" {
  id = tencentcloud_postgresql_instance.example.id
}

data "tencentcloud_postgresql_instances" "project_example" {
  project_id = 0
}

data "tencentcloud_postgresql_instances" "name_example" {
  name = tencentcloud_postgresql_instance.example.name
}

data "tencentcloud_postgresql_specinfos" "example" {
  availability_zone = var.availability_zone
}
