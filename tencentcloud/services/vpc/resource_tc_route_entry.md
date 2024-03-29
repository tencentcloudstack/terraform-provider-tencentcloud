Provides a resource to create a routing entry in a VPC routing table.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_route_table_entry.

Example Usage

```hcl
resource "tencentcloud_vpc" "main" {
  name       = "Used to test the routing entry"
  cidr_block = "10.4.0.0/16"
}

resource "tencentcloud_route_table" "r" {
  name   = "Used to test the routing entry"
  vpc_id = tencentcloud_vpc.main.id
}

resource "tencentcloud_route_entry" "rtb_entry_instance" {
  vpc_id         = tencentcloud_route_table.main.vpc_id
  route_table_id = tencentcloud_route_table.r.id
  cidr_block     = "10.4.8.0/24"
  next_type      = "instance"
  next_hub       = "10.16.1.7"
}

resource "tencentcloud_route_entry" "rtb_entry_instance" {
  vpc_id         = tencentcloud_route_table.main.vpc_id
  route_table_id = tencentcloud_route_table.r.id
  cidr_block     = "10.4.5.0/24"
  next_type      = "vpn_gateway"
  next_hub       = "vpngw-db52irtl"
}
```