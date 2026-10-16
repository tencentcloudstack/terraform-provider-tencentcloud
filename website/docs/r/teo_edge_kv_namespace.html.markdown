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

* Create: [CreateEdgeKVNamespace](https://www.tencentcloud.com/search/?wd=CreateEdgeKVNamespace&type=document&lang=en)
* Create: [DescribeEdgeKVNamespaces](https://www.tencentcloud.com/search/?wd=DescribeEdgeKVNamespaces&type=document&lang=en)
* Read: [DescribeEdgeKVNamespaces](https://www.tencentcloud.com/search/?wd=DescribeEdgeKVNamespaces&type=document&lang=en)
* Update: [DescribeEdgeKVNamespaces](https://www.tencentcloud.com/search/?wd=DescribeEdgeKVNamespaces&type=document&lang=en)
* Update: [ModifyEdgeKVNamespace](https://www.tencentcloud.com/search/?wd=ModifyEdgeKVNamespace&type=document&lang=en)
* Delete: [DeleteEdgeKVNamespace](https://www.tencentcloud.com/search/?wd=DeleteEdgeKVNamespace&type=document&lang=en)

## Import

TEO Edge KV namespace can be imported using the zone_id#namespace, e.g.

```
terraform import tencentcloud_teo_edge_kv_namespace.example zone-2o3h21ed2t68#example-namespace
```

