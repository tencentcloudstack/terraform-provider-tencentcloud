---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_multi_path_gateway"
sidebar_current: "docs-tencentcloud-resource-teo_multi_path_gateway"
description: |-
  Provides a resource to create a teo multi path gateway for EdgeOne(TEO).
---

# tencentcloud_teo_multi_path_gateway

Provides a resource to create a teo multi path gateway for EdgeOne(TEO).

## Example Usage

### Cloud type gateway

```hcl
resource "tencentcloud_teo_multi_path_gateway" "cloud" {
  zone_id      = "zone-359h792djt7h"
  gateway_type = "cloud"
  gateway_name = "test-cloud-gw"
  region_id    = "ap-guangzhou"
}
```

### Private type gateway

```hcl
resource "tencentcloud_teo_multi_path_gateway" "private" {
  zone_id      = "zone-359h792djt7h"
  gateway_type = "private"
  gateway_name = "test-private-gw"
  gateway_ip   = "1.2.3.4"
  gateway_port = 8080
}
```

## Argument Reference

The following arguments are supported:

* `gateway_name` - (Required, String) Gateway name, up to 16 characters.
* `gateway_type` - (Required, String, ForceNew) Gateway type. Valid values: `cloud`, `private`.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `gateway_ip` - (Optional, String) Gateway IP address, required when GatewayType is private.
* `gateway_port` - (Optional, Int) Gateway port, range 1-65535 (excluding 8888).
* `region_id` - (Optional, String, ForceNew) Gateway region, required when GatewayType is cloud.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `gateway_id` - Gateway ID.
* `need_confirm` - Whether the gateway origin IP list needs reconfirmation.
* `status` - Gateway status.


## Import

teo multi path gateway can be imported using the id, e.g.

```
terraform import tencentcloud_teo_multi_path_gateway.example zone-279qso5a4cw9#mpgw-g3176ppeye
```

