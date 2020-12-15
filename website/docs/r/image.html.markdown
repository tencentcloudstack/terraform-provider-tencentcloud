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

## Argument Reference

The following arguments are supported:

* `image_name` - (Required) Image name.
* `data_disk_ids` - (Optional, ForceNew) Cloud disk ID list, When creating a whole machine image based on an instance, specify the data disk ID contained in the image.
* `force_poweroff` - (Optional) Set whether to force shutdown during mirroring. The default value is `false`, when set to true, it means that the mirror will be made after shutdown.
* `image_description` - (Optional) Image Description.
* `instance_id` - (Optional, ForceNew) Cloud server instance ID.
* `snapshot_ids` - (Optional, ForceNew) Cloud disk snapshot ID list; creating a mirror based on a snapshot must include a system disk snapshot. It cannot be passed in simultaneously with InstanceId.
* `sysprep` - (Optional) Sysprep function under Windows. When creating a Windows image, you can select true or false to enable or disable the Syspre function.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

image instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_image.image_snap img-gf7jspk6
```

