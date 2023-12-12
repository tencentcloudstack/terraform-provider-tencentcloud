Provides details about a specific Route Table.

This resource can prove useful when a module accepts a Subnet id as an input variable and needs to, for example, add a route in the Route Table.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_vpc_route_tables.

Example Usage

```hcl
variable "route_table_id" {}

data "tencentcloud_route_table" "selected" {
  route_table_id = var.route_table_id
}

resource "tencentcloud_route_entry" "rtb_entry_instance" {
  vpc_id         = "{data.tencentcloud_route_table.selected.vpc_id}"
  route_table_id = var.route_table_id
  cidr_block     = "10.4.8.0/24"
  next_type      = "instance"
  next_hub       = "10.16.1.7"
}
```