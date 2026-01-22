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
data "tencentcloud_images" "example" {
  image_type       = ["PRIVATE_IMAGE"]
  image_name_regex = "MyImage"
}

resource "tencentcloud_cvm_sync_image" "example" {
  image_id            = data.tencentcloud_images.example.images.0.image_id
  destination_regions = ["ap-guangzhou", "ap-shanghai"]
  encrypt             = true
  kms_key_id          = "f063c18b-654b-11ef-9d9f-525400d3a886"
}
```

## Argument Reference

The following arguments are supported:

* `destination_regions` - (Required, Set: [`String`], ForceNew) List of destination regions for synchronization. Limits: It must be a valid region. For a custom image, the destination region cannot be the source region. For a shared image, the destination region must be the source region, which indicates to create a copy of the image as a custom image in the same region.
* `image_id` - (Required, String, ForceNew) Image ID. The specified image must meet the following requirement: the images must be in the `NORMAL` state.
* `dry_run` - (Optional, Bool, ForceNew) Checks whether image synchronization can be initiated.
* `encrypt` - (Optional, Bool, ForceNew) Whether to synchronize as an encrypted custom image. Default value is `false`. Synchronization to an encrypted custom image is only supported within the same region.
* `image_name` - (Optional, String, ForceNew) Destination image name.
* `image_set_required` - (Optional, Bool, ForceNew) Whether to return the ID of image created in the destination region.
* `kms_key_id` - (Optional, String, ForceNew) KMS key ID used when synchronizing to an encrypted custom image. This parameter is valid only synchronizing to an encrypted image. If KmsKeyId is not specified, the default CBS cloud product KMS key is used.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `image_set` - ID of the image created in the destination region.
  * `image_id` - Image ID.
  * `region` - Region of the image.


