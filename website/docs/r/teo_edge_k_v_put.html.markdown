---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_edge_k_v_put"
sidebar_current: "docs-tencentcloud-resource-teo_edge_k_v_put"
description: |-
  Provides a resource to create a TEO Edge KV key-value pair attachment
---

# tencentcloud_teo_edge_k_v_put

Provides a resource to create a TEO Edge KV key-value pair attachment

## Example Usage

```hcl
resource "tencentcloud_teo_edge_k_v_put" "example" {
  zone_id   = "zone-2o3h21ed2t68"
  namespace = "example-namespace"
  key       = "example-key"
  value     = "example-value"
}
```

### With expiration

```hcl
resource "tencentcloud_teo_edge_k_v_put" "example_with_ttl" {
  zone_id        = "zone-2o3h21ed2t68"
  namespace      = "example-namespace"
  key            = "example-key-ttl"
  value          = "example-value"
  expiration_ttl = 3600
}
```

## Argument Reference

The following arguments are supported:

* `key` - (Required, String, ForceNew) Key name, 1-512 characters, allowed characters are letters, numbers, hyphens and underscores.
* `namespace` - (Required, String, ForceNew) Namespace name.
* `value` - (Required, String) Key value. Cannot be empty, maximum 1 MB.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `expiration_ttl` - (Optional, Int) Time-to-live of the key-value pair, relative time in seconds. Must be greater than or equal to 60.
* `expiration` - (Optional, Int) Expiration time of the key-value pair, absolute time, Unix timestamp in seconds. Must be greater than or equal to current time + 60.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

TEO Edge KV put can be imported using the zoneId#namespace#key, e.g.

```
terraform import tencentcloud_teo_edge_k_v_put.example zone-2o3h21ed2t68#example-namespace#example-key
```

