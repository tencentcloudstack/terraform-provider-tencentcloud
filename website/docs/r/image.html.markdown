---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_image"
sidebar_current: "docs-tencentcloud-resource-image"
description: |-
  Provide a resource to manage image.
---

# tencentcloud_image

Provide a resource to manage image.

## Example Usage

```hcl
resource "tencentcloud_image" "image_snap" {
  image_name        = "image-snapshot-keep"
  snapshot_ids      = ["snap-nbp3xy1d", "snap-nvzu3dmh"]
  force_poweroff    = true
  image_description = "create image with snapshot"
}
```

### Use image family

```hcl
resource "tencentcloud_image" "image_family" {
  image_description = "create image with snapshot 12323"
  image_family      = "business-daily-update"
  image_name        = "image-family-test123"
  data_disk_ids     = []
  snapshot_ids = [
    "snap-7uuvrcoj",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `image_name` - (Required, String) Image name.
* `data_disk_ids` - (Optional, Set: [`String`], ForceNew) Cloud disk ID list, When creating a whole machine image based on an instance, specify the data disk ID contained in the image.
* `force_poweroff` - (Optional, Bool) Set whether to force shutdown during mirroring. The default value is `false`, when set to true, it means that the mirror will be made after shutdown.
* `image_description` - (Optional, String) Image Description.
* `image_family` - (Optional, String) Set image family. Example value: `business-daily-update`.
* `instance_id` - (Optional, String, ForceNew) Cloud server instance ID.
* `snapshot_ids` - (Optional, Set: [`String`], ForceNew) Cloud disk snapshot ID list; creating a mirror based on a snapshot must include a system disk snapshot. It cannot be passed in simultaneously with InstanceId.
* `sysprep` - (Optional, Bool) Sysprep function under Windows. When creating a Windows image, you can select true or false to enable or disable the Syspre function.
* `tags` - (Optional, Map) Tags of the image.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateImage](https://www.tencentcloud.com/document/product/213/33276) - This API is used to create a new image with the system disk of an instance. The image can be used to create new instances.
* Create: [DescribeImages](https://www.tencentcloud.com/document/product/213/33272) - This API is used to view the list of images.
* Create: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Read: [DescribeImages](https://www.tencentcloud.com/document/product/213/33272) - This API is used to view the list of images.
* Read: [DescribeResourceTagsByResourceIds](https://www.tencentcloud.com/document/product/651/33022) - This API is used to query the tag key-value pairs associated with an existing resource.
* Update: [ModifyImageAttribute](https://www.tencentcloud.com/document/product/213/33269) - This API is used to modify image attributes.
* Update: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Delete: [DeleteImages](https://www.tencentcloud.com/document/product/213/33275) - This API is used to delete one or more images.
* Delete: [DescribeImages](https://www.tencentcloud.com/document/product/213/33272) - This API is used to view the list of images.

## Import

image instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_image.image_snap img-gf7jspk6
```

