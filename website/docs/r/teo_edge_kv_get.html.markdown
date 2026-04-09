---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_edge_kv_get"
sidebar_current: "docs-tencentcloud-resource-teo_edge_kv_get"
description: |-
  Provides a resource to query Edge KV data
---

# tencentcloud_teo_edge_kv_get

Provides a resource to query Edge KV (Edge Key-Value) data from TencentCloud EdgeOne service.

~> **NOTE:** This is a query-only resource. It retrieves existing KV data from the EdgeOne service without creating or modifying any data.

## Example Usage

### Basic Usage

```hcl
resource "tencentcloud_teo_edge_kv_get" "example" {
  zone_id   = "zone-3j1xw7910arp"
  namespace = "ns-011"
  keys      = ["hello", "world"]
}
```

### Query Multiple Keys

```hcl
resource "tencentcloud_teo_edge_kv_get" "multi_keys" {
  zone_id   = "zone-3j1xw7910arp"
  namespace = "ns-011"
  keys      = ["key1", "key2", "key3"]
}
```

### Query with Namespace

```hcl
data "tencentcloud_teo_edge_kv_namespaces" "example" {
  zone_id = "zone-3j1xw7910arp"
}

resource "tencentcloud_teo_edge_kv_get" "example" {
  zone_id   = "zone-3j1xw7910arp"
  namespace = data.tencentcloud_teo_edge_kv_namespaces.example.namespaces[0].name
  keys      = ["config_key"]
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) Zone ID. Example: `zone-3j1xw7910arp`
* `namespace` - (Required, ForceNew) Namespace name. You can get the list of namespaces under a site through the `DescribeEdgeKVNamespaces` interface. Example: `ns-011`
* `keys` - (Required) List of key names to query.
  * Each key name must be between 1-512 characters
  * Key names can only contain letters, numbers, hyphens, and underscores
  * The list must contain at least 1 element
  * The list must contain at most 20 elements
  * Example: `["hello", "world"]`

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in the format of `zoneId#namespace#firstKey`.
* `data` - List of key-value pair data, returned in the same order as the `keys` parameter. If a key does not exist, the `value` field for that key returns an empty string.
  * `key` - The key name.
  * `value` - The key value. If the key does not exist, an empty string is returned.
  * `expiration` - Expiration time in ISO 8601 format (YYYY-MM-DDThh:mm:ssZ, UTC time). If empty, the key-value pair never expires.

## Import

This resource does not support import as it is a query-only resource.

## Notes

* This resource retrieves data from the EdgeOne Edge KV service.
* The `data` field contains the key-value pairs in the same order as the `keys` parameter.
* If a queried key does not exist in the namespace, its `value` field will return an empty string.
* The `expiration` field is an empty string if the key-value pair does not have an expiration time.
* Deleting this resource only removes it from the Terraform state and does not affect the actual KV data in the EdgeOne service.
* The resource ID is constructed as `zoneId#namespace#firstKey` where `firstKey` is the first element of the `keys` list.
* Changing `zone_id` or `namespace` will force recreation of the resource.
