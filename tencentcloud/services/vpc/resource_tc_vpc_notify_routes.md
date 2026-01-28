Provides a resource to create a VPC notify routes

Example Usage

```hcl
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_route_table" "route_table" {
  vpc_id = tencentcloud_vpc.vpc.id
  name   = "tf-example"
}

resource "tencentcloud_vpc_notify_routes" "example" {
  route_table_id = tencentcloud_route_table.route_table.id
  route_item_ids = ["rti-i8bap903"]
}
```

Import

VPC notify routes can be imported using the routeTableId#routeItemId, e.g.

```
terraform import tencentcloud_vpc_notify_routes.example route_table_id#route_item_id
```
