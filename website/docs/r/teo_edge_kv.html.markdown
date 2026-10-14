---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_edge_kv"
sidebar_current: "docs-tencentcloud-resource-teo_edge_kv"
description: |-
  Provides a resource to create a TEO Edge KV key-value pair
---

# tencentcloud_teo_edge_kv

Provides a resource to create a TEO Edge KV key-value pair

## Example Usage

```hcl
resource "tencentcloud_teo_edge_kv" "example" {
  zone_id   = "zone-2o3h21ed2t68"
  namespace = "example-namespace"
  key       = "example-key"
  value     = "example-value"
}
```

## Argument Reference

The following arguments are supported:

* `key` - (Required, String, ForceNew) Key name, 1-512 characters, allowed characters are letters, numbers, hyphens and underscores.
* `namespace` - (Required, String, ForceNew) Namespace name.
* `value` - (Required, String) Key value. Cannot be empty, maximum 1 MB.
* `zone_id` - (Required, String, ForceNew) Site ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [EdgeKVGet](https://cloud.tencent.com/document/api/1552/130185) - This API is used to batch read key values from a specified namespace, supporting a single query of up to 20 keys.
* Create: [EdgeKVPut](https://cloud.tencent.com/document/api/1552/130183) - This API is used to write key-value pair data to a specified namespace and supports setting expiration time. If the key already exists, it overwrites the original value. If Not Exist, it creates a new key-value pair.
* Read: [EdgeKVGet](https://cloud.tencent.com/document/api/1552/130185) - This API is used to batch read key values from a specified namespace, supporting a single query of up to 20 keys.
* Update: [EdgeKVGet](https://cloud.tencent.com/document/api/1552/130185) - This API is used to batch read key values from a specified namespace, supporting a single query of up to 20 keys.
* Update: [EdgeKVPut](https://cloud.tencent.com/document/api/1552/130183) - This API is used to write key-value pair data to a specified namespace and supports setting expiration time. If the key already exists, it overwrites the original value. If Not Exist, it creates a new key-value pair.
* Delete: [EdgeKVDelete](https://cloud.tencent.com/document/api/1552/130186) - This API is used to delete one or more key-value pair data in the specified namespace, with batch deletion supported. Data cannot be recovered after deletion.

## Import

TEO Edge KV can be imported using the zoneId#namespace#key, e.g.

```
terraform import tencentcloud_teo_edge_kv.example zone-2o3h21ed2t68#example-namespace#example-key
```

