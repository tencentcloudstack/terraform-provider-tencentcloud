Provides a resource to create a Route table entry.

~> **NOTE:** When setting the route item switch, do not use it together with resource `tencentcloud_route_table_entry_config`.

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create route table
resource "tencentcloud_route_table" "example" {
  vpc_id = tencentcloud_vpc.vpc.id
  name   = "tf-example"
}

# create subnet
resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet"
  cidr_block        = "10.0.12.0/24"
  availability_zone = var.availability_zone
  route_table_id    = tencentcloud_route_table.example.id
}

# create route table entry
resource "tencentcloud_route_table_entry" "example" {
  route_table_id         = tencentcloud_route_table.example.id
  destination_cidr_block = "10.12.12.0/24"
  next_type              = "EIP"
  next_hub               = "0"
  description            = "Terraform test."
}

# output
output "item_id" {
  value = tencentcloud_route_table_entry.example.route_item_id
}
```

Import

Route table entry can be imported using the routeEntryId.routeTableId, e.g.

```
$ terraform import tencentcloud_route_table_entry.example 3065857.rtb-b050fg94
```
