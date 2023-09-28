---
subcategory: "Cfw"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_sync_route"
sidebar_current: "docs-tencentcloud-resource-cfw_sync_route"
description: |-
  Provides a resource to create a cfw sync_route
---

# tencentcloud_cfw_sync_route

Provides a resource to create a cfw sync_route

## Example Usage

```hcl
resource "tencentcloud_cfw_sync_route" "example" {
  sync_type = "Route"
  fw_type   = "nat"
}
```

## Argument Reference

The following arguments are supported:

* `sync_type` - (Required, String, ForceNew) Synchronization operation type: Route, synchronize firewall routing.
* `fw_type` - (Optional, String, ForceNew) Firewall type; nat: nat firewall; ew: inter-vpc firewall.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



