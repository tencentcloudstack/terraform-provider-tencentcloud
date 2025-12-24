Provides a resource to create a VPC route policy

Example Usage

```hcl
resource "tencentcloud_vpc_route_policy" "example" {
  route_policy_name        = "tf-example"
  route_policy_description = "remark."
}
```

Import

VPC route policy can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_route_policy.example rrp-lpv8rjp8
```
