resource "tencentcloud_vpc" "main" {
  name       = "${var.short_name}"
  cidr_block = "${var.vpc_cidr}"
}
