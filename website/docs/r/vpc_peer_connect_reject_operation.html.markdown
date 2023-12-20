---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_peer_connect_reject_operation"
sidebar_current: "docs-tencentcloud-resource-vpc_peer_connect_reject_operation"
description: |-
  Provides a resource to create a vpc peer_connect_reject_operation
---

# tencentcloud_vpc_peer_connect_reject_operation

Provides a resource to create a vpc peer_connect_reject_operation

## Example Usage

```hcl
resource "tencentcloud_vpc_peer_connect_reject_operation" "peer_connect_reject_operation" {
  peering_connection_id = "pcx-abced"
}
```

## Argument Reference

The following arguments are supported:

* `peering_connection_id` - (Required, String, ForceNew) Peer connection unique ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc peer_connect_reject_operation can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_peer_connect_reject_operation.peer_connect_reject_operation peer_connect_reject_operation_id
```

