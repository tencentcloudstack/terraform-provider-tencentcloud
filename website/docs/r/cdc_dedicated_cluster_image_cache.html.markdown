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


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateDedicatedClusterImageCache](https://www.tencentcloud.com/search/?wd=CreateDedicatedClusterImageCache&type=document&lang=en)
  - [DescribeImages](https://www.tencentcloud.com/document/api/213/33272)
* Read: [DescribeImages](https://www.tencentcloud.com/document/api/213/33272)
* Delete:
  - [DeleteDedicatedClusterImageCache](https://www.tencentcloud.com/search/?wd=DeleteDedicatedClusterImageCache&type=document&lang=en)
  - [DescribeImages](https://www.tencentcloud.com/document/api/213/33272)

## Import

CDC dedicated cluster image cache can be imported using the id, e.g.

```
terraform import tencentcloud_cdc_dedicated_cluster_image_cache.cdc_dedicated_cluster_image_cache ${dedicated_cluster_id}#${image_id}
```

