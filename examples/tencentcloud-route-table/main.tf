resource "tencentcloud_vpc" "my_vpc" {
  name       = "Used to test rtb"
  cidr_block = "10.2.0.0/16"
}

resource "tencentcloud_route_table" "my_rtb" {
  vpc_id = "${tencentcloud_vpc.my_vpc.id}"
  name   = "${var.short_name}"

  tags = {
    "test" = "test"
  }
}

data "tencentcloud_vpc_route_tables" "tags_instances" {
  tags = "${tencentcloud_route_table.my_rtb.tags}"
}
