resource "tencentcloud_vpc" "my_vpc" {
  name       = "${var.short_name}"
  cidr_block = "${var.vpc_cidr}"
}

resource "tencentcloud_route_table" "my_rtb" {
  vpc_id = "${tencentcloud_vpc.my_vpc.id}"
  name   = "Used to test rtb entry"
}

resource "tencentcloud_route_entry" "rtb_entry_nat" {
  vpc_id         = "${tencentcloud_route_table.my_rtb.vpc_id}"
  route_table_id = "${tencentcloud_route_table.my_rtb.id}"
  cidr_block     = "10.4.4.0/24"
  next_type      = "nat_gateway"
  next_hub       = "nat-dt9ycr9y" //note this is hardcode, need to replace it in your real situation
}

resource "tencentcloud_route_entry" "rtb_entry_vpn" {
  vpc_id         = "${tencentcloud_route_table.my_rtb.vpc_id}"
  route_table_id = "${tencentcloud_route_table.my_rtb.id}"
  cidr_block     = "10.4.5.0/24"
  next_type      = "vpn_gateway"
  next_hub       = "vpngw-db52irtl" //note this is hardcode, need to replace it in your real situation
}

resource "tencentcloud_route_entry" "rtb_entry_dc" {
  vpc_id         = "${tencentcloud_route_table.my_rtb.vpc_id}"
  route_table_id = "${tencentcloud_route_table.my_rtb.id}"
  cidr_block     = "10.4.6.0/24"
  next_type      = "dc_gateway"
  next_hub       = "dcg-9r7vi45r" //note this is hardcode, need to replace it in your real situation
}

resource "tencentcloud_route_entry" "rtb_entry_peering_connection" {
  vpc_id         = "${tencentcloud_route_table.my_rtb.vpc_id}"
  route_table_id = "${tencentcloud_route_table.my_rtb.id}"
  cidr_block     = "172.17.2.0/24"
  next_type      = "peering_connection"
  next_hub       = "pcx-bj5b69qu" //note this is hardcode, need to replace it in your real situation
}

resource "tencentcloud_route_entry" "rtb_entry_instance" {
  vpc_id         = "${tencentcloud_route_table.my_rtb.vpc_id}"
  route_table_id = "${tencentcloud_route_table.my_rtb.id}"
  cidr_block     = "10.4.8.0/24"
  next_type      = "instance"
  next_hub       = "10.16.1.7" //note this is hardcode, need to replace it in your real situation
}
