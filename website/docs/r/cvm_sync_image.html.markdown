---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_sync_image"
sidebar_current: "docs-tencentcloud-resource-cvm_sync_image"
description: |-
  Provides a resource to create a cvm sync_image
---

# tencentcloud_cvm_sync_image

Provides a resource to create a cvm sync_image

## Example Usage

```hcl
resource "tencentcloud_cvm_sync_image" "sync_image" {
  image_id            = "img-xxxxxx"
  destination_regions = ["ap-guangzhou", "ap-shanghai"]
}
```

## Argument Reference

The following arguments are supported:

* `destination_regions` - (Required, Set: [`String`], ForceNew) List of destination regions for synchronization. Limits: It must be a valid region. For a custom image, the destination region cannot be the source region. For a shared image, the destination region must be the source region, which indicates to create a copy of the image as a custom image in the same region.
* `image_id` - (Required, String, ForceNew) Image ID. The specified image must meet the following requirement: the images must be in the `NORMAL` state.
* `dry_run` - (Optional, Bool, ForceNew) Checks whether image synchronization can be initiated.
* `image_name` - (Optional, String, ForceNew) Destination image name.
* `image_set_required` - (Optional, Bool, ForceNew) Whether to return the ID of image created in the destination region.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



