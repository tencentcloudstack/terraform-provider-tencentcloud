---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_multi_path_gateway_line"
sidebar_current: "docs-tencentcloud-resource-teo_multi_path_gateway_line"
description: |-
  Provides a resource to create a teo multi path gateway line for EdgeOne(TEO).
---

# tencentcloud_teo_multi_path_gateway_line

Provides a resource to create a teo multi path gateway line for EdgeOne(TEO).

## Example Usage

### Custom line type

```hcl
resource "tencentcloud_teo_multi_path_gateway_line" "example" {
  zone_id      = "zone-279qso5a4cw9"
  gateway_id   = "gw-2qwk1t3g3jx9"
  line_type    = "custom"
  line_address = "1.2.3.4:80"
}
```

### Proxy line type

```hcl
resource "tencentcloud_teo_multi_path_gateway_line" "example" {
  zone_id      = "zone-279qso5a4cw9"
  gateway_id   = "gw-2qwk1t3g3jx9"
  line_type    = "proxy"
  line_address = "5.6.7.8:443"
  proxy_id     = "sid-38hbn26osico"
  rule_id      = "rule-abcdef"
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String, ForceNew) Multi-path gateway ID.
* `line_address` - (Required, String) Line address, format is host:port.
* `line_type` - (Required, String) Line type. Valid values: `direct`, `proxy`, `custom`.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `proxy_id` - (Optional, String) L4 proxy instance ID, required when LineType is proxy.
* `rule_id` - (Optional, String) Forwarding rule ID, required when LineType is proxy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `line_id` - Line ID, returned after creation by cloud API.


## Import

teo multi path gateway line can be imported using the id, e.g.

```
terraform import tencentcloud_teo_multi_path_gateway_line.example zone-279qso5a4cw9#gw-2qwk1t3g3jx9#line-1
```

