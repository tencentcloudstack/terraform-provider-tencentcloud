---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_route_policy_entries"
sidebar_current: "docs-tencentcloud-resource-vpc_route_policy_entries"
description: |-
  Provides a resource to create a VPC route policy entries
---

# tencentcloud_vpc_route_policy_entries

Provides a resource to create a VPC route policy entries

~> **NOTE:** This resource must exclusive in one route policy ID, do not declare additional route policy entries resources of this route policy ID elsewhere.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `route_policy_entry_set` - (Required, Set) Route reception policy entry list.
* `route_policy_id` - (Required, String, ForceNew) Specifies the instance ID of the route reception policy.

The `route_policy_entry_set` object supports the following:

* `action` - (Optional, String) Action.
DROP: drop.
DISABLE: receive and disable.
ACCEPT: receive and enable.
Note: This field may return null, indicating that no valid value was found.
* `cidr_block` - (Optional, String) Destination ip range.
Note: This field may return null, indicating that no valid value was found.
* `description` - (Optional, String) Describes the routing strategy rule.
Note: This field may return null, indicating that no valid value was found.
* `gateway_id` - (Optional, String) Gateway unique ID.
Note: This field may return null, indicating that no valid value was found.
* `gateway_type` - (Optional, String) Next hop type. types currently supported:.
CVM: cloud virtual machine with public network gateway type.
VPN: vpn gateway.
DIRECTCONNECT: direct connect gateway.
PEERCONNECTION: peering connection.
HAVIP: high availability virtual ip.
NAT: specifies the nat gateway. 
EIP: specifies the public ip address of the cloud virtual machine.
LOCAL_GATEWAY: specifies the local gateway.
PVGW: pvgw gateway.
Note: This field may return null, indicating that no valid value was found.
* `priority` - (Optional, Int) Priority. a smaller value indicates a higher priority.
Note: This field may return null, indicating that no valid value was found.
* `route_type` - (Optional, String) Routing Type

Specifies the USER-customized data type.
NETD: specifies the route for network detection.
CCN: CCN route.
Note: This field may return null, indicating that no valid value was found.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

VPC route policy entries can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_route_policy_entries.example rrp-lpv8rjp8
```

