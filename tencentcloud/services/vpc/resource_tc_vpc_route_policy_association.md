Provides a resource to create a VPC route policy association

Example Usage

```hcl
resource "tencentcloud_vpc_route_policy_association" "example" {
  route_policy_id = "rrp-7dnu4yoi"
  route_table_id  = "rtb-389phpuq"
  priority        = 10
}
```

Import

VPC route policy association can be imported using the routePolicyId#routeTableId, e.g.

```
terraform import tencentcloud_vpc_route_policy_association.example rrp-7dnu4yoi#rtb-389phpuq
```
