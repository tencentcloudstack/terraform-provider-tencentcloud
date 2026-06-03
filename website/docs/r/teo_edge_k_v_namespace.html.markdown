---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_edge_k_v_namespace"
sidebar_current: "docs-tencentcloud-resource-teo_edge_k_v_namespace"
description: |-
  Provides a resource to create a TEO edge KV namespace
---

# tencentcloud_teo_edge_k_v_namespace

Provides a resource to create a TEO edge KV namespace

## Example Usage

```hcl
resource "tencentcloud_teo_edge_k_v_namespace" "example" {
  zone_id   = "zone-2o3h21ed2t68"
  namespace = "my-namespace"
  remark    = "This is an example namespace"
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Required, String, ForceNew) Namespace name. Supports 1-50 characters, allowed characters are a-z, A-Z, 0-9, -.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `remark` - (Optional, String) Namespace description. Maximum 256 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

TEO edge KV namespace can be imported using the composite id (zone_id#namespace), e.g.

````
terraform import tencentcloud_teo_edge_k_v_namespace.example zone-2o3h21ed2t68#my-namespace
````

