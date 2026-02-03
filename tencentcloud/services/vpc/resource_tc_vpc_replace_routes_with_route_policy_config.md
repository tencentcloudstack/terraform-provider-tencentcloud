Provides a resource to create a VPC replace routes with route policy config

Example Usage

```hcl
resource "tencentcloud_vpc_replace_routes_with_route_policy_config" "example" {
  route_table_id = "rtb-olsbhnyc"
  routes {
    route_item_id      = "rti-araogi5t"
    force_match_policy = true
  }

  routes {
    route_item_id      = "rti-kiyt72op"
    force_match_policy = true
  }
}
```
