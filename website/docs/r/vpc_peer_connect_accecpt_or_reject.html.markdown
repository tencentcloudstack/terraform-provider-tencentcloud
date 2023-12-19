---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_peer_connect_accecpt_or_reject"
sidebar_current: "docs-tencentcloud-resource-vpc_peer_connect_accecpt_or_reject"
description: |-
  Provides a resource to create a vpc peer_connect_accecpt_or_reject
---

# tencentcloud_vpc_peer_connect_accecpt_or_reject

Provides a resource to create a vpc peer_connect_accecpt_or_reject

## Example Usage

```hcl
resource "tencentcloud_vpc_peer_connect_accecpt_or_reject" "peer_connect_accecpt_or_reject" {
  peering_connection_id = "pcx-abced"
}
```

## Argument Reference

The following arguments are supported:

* `peering_connection_id` - (Required, String, ForceNew) Peer connection unique ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



