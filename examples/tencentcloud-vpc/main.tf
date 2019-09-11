resource "tencentcloud_vpc" "main" {
  name       = "${var.short_name}"
  cidr_block = "${var.vpc_cidr}"

  tags = {
    "test" = "test"
  }
}

data "tencentcloud_vpc_instances" "tags_instances" {
  name = "${tencentcloud_vpc.main.name}"
  tags = "${tencentcloud_vpc.main.tags}"
}