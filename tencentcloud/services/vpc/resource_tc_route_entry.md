Provides a resource to create a routing entry in a VPC routing table.

~> **NOTE:** It has been deprecated and replaced by `tencentcloud_route_table_entry`.

Example Usage

```hcl
resource "tencentcloud_vpc" "vpc" {
  name       = "tf-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_route_table" "example" {
  name   = "tf-example"
  vpc_id = tencentcloud_vpc.vpc.id
}

resource "tencentcloud_route_entry" "example1" {
  vpc_id         = tencentcloud_vpc.vpc.id
  route_table_id = tencentcloud_route_table.example.id
  cidr_block     = "192.168.0.0/24"
  next_type      = "eip"
  next_hub       = "0"
}

resource "tencentcloud_route_entry" "example2" {
  vpc_id         = tencentcloud_vpc.vpc.id
  route_table_id = tencentcloud_route_table.example.id
  cidr_block     = "10.4.8.0/24"
  next_type      = "instance"
  next_hub       = "10.16.1.7"
}
```