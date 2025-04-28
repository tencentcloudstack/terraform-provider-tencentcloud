---
subcategory: "CDC"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdc_dedicated_cluster_image_cache"
sidebar_current: "docs-tencentcloud-resource-cdc_dedicated_cluster_image_cache"
description: |-
  Provides a resource to create a CDC dedicated cluster image cache
---

# tencentcloud_cdc_dedicated_cluster_image_cache

Provides a resource to create a CDC dedicated cluster image cache

## Example Usage

```hcl
resource "tencentcloud_cdc_dedicated_cluster_image_cache" "cdc_dedicated_cluster_image_cache" {
  dedicated_cluster_id = "cluster-262n63e8"
  image_id             = "img-eb30mz89"
}
```

## Argument Reference

The following arguments are supported:

* `dedicated_cluster_id` - (Required, String, ForceNew) Cluster ID.
* `image_id` - (Required, String, ForceNew) Image ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CDC dedicated cluster image cache can be imported using the id, e.g.

```
terraform import tencentcloud_cdc_dedicated_cluster_image_cache.cdc_dedicated_cluster_image_cache ${dedicated_cluster_id}#${image_id}
```

