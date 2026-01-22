Provides a resource to create a VPC route policy entries

~> **NOTE:** This resource must exclusive in one route policy ID, do not declare additional route policy entries resources of this route policy ID elsewhere.

Example Usage

```hcl
resource "tencentcloud_vpc_route_policy_entries" "example" {
  route_policy_id = tencentcloud_vpc_route_policy.example.id
  route_policy_entry_set {
    cidr_block   = "10.10.0.0/16"
    route_type   = "ANY"
    gateway_type = "VPN"
    gateway_id   = "vpngw-may3cb0m"
    action       = "ACCEPT"
  }

  route_policy_entry_set {
    cidr_block   = "172.16.0.0/16"
    description  = "remark"
    route_type   = "ANY"
    gateway_type = "EIP"
    priority     = 10
    action       = "ACCEPT"
  }

  route_policy_entry_set {
    cidr_block   = "192.168.0.0/16"
    description  = "remark"
    route_type   = "ANY"
    gateway_type = "HAVIP"
    gateway_id   = "havip-r3ar5p86"
    priority     = 1
    action       = "ACCEPT"
  }
}
```

Import

VPC route policy entries can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_route_policy_entries.example rrp-lpv8rjp8
```
