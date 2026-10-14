---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_edge_kv_namespace"
sidebar_current: "docs-tencentcloud-resource-teo_edge_kv_namespace"
description: |-
  Provides a resource to create a TEO Edge KV namespace
---

# tencentcloud_teo_edge_kv_namespace

Provides a resource to create a TEO Edge KV namespace

## Example Usage

```hcl
resource "tencentcloud_teo_edge_kv_namespace" "example" {
  zone_id   = "zone-2o3h21ed2t68"
  namespace = "example-namespace"
  remark    = "This is an example namespace"
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Required, String, ForceNew) Namespace name. Supports 1-50 characters, allowed characters are a-z, A-Z, 0-9, -, and - cannot be used alone or consecutively, cannot be placed at the beginning or end. The name must be unique within the same site.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `remark` - (Optional, String) Namespace description. Used to describe the purpose or business meaning of the namespace. Maximum 256 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `capacity_used` - KV storage used capacity in bytes.
* `capacity` - KV storage available capacity in bytes.
* `created_on` - Creation time in ISO 8601 format (UTC).
* `modified_on` - Last modification time in ISO 8601 format (UTC).

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateEdgeKVNamespace](https://cloud.tencent.com/document/api/1552/130189) - This API is used to create a KV namespace in the specified site.
* Create: [DescribeEdgeKVNamespaces](https://cloud.tencent.com/document/api/1552/130187) - This API is used to query the KV namespace list of a specified site, supporting pagination, sorting and conditional filtering. It returns the basic info, capacity utilization and reference relationship of namespaces.
* Read: [DescribeEdgeKVNamespaces](https://cloud.tencent.com/document/api/1552/130187) - This API is used to query the KV namespace list of a specified site, supporting pagination, sorting and conditional filtering. It returns the basic info, capacity utilization and reference relationship of namespaces.
* Update: [DescribeEdgeKVNamespaces](https://cloud.tencent.com/document/api/1552/130187) - This API is used to query the KV namespace list of a specified site, supporting pagination, sorting and conditional filtering. It returns the basic info, capacity utilization and reference relationship of namespaces.
* Update: [ModifyEdgeKVNamespace](https://cloud.tencent.com/document/api/1552/130182) - This API is used to modify the attribute information of a specified KV namespace. Currently supported is namespace description modification.
* Delete: [DeleteEdgeKVNamespace](https://cloud.tencent.com/document/api/1552/130188) - This API is used to delete specified KV namespace. Once deleted, all keys in namespace will be cleared and cannot be recovered.

## Import

TEO Edge KV namespace can be imported using the zone_id#namespace, e.g.

```
terraform import tencentcloud_teo_edge_kv_namespace.example zone-2o3h21ed2t68#example-namespace
```

