---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_export_images"
sidebar_current: "docs-tencentcloud-resource-cvm_export_images"
description: |-
  Provides a resource to create a cvm export_images
---

# tencentcloud_cvm_export_images

Provides a resource to create a cvm export_images

## Example Usage

```hcl
resource "tencentcloud_cvm_export_images" "export_images" {
  bucket_name      = "xxxxxx"
  image_id         = "img-xxxxxx"
  file_name_prefix = "test-"
}
```

## Argument Reference

The following arguments are supported:

* `bucket_name` - (Required, String, ForceNew) COS bucket name.
* `file_name_prefix` - (Required, String, ForceNew) Prefix of exported file.
* `image_id` - (Required, String, ForceNew) Image ID.
* `dry_run` - (Optional, Bool, ForceNew) Check whether the image can be exported.
* `export_format` - (Optional, String, ForceNew) Format of the exported image file. Valid values: RAW, QCOW2, VHD and VMDK. Default value: RAW.
* `only_export_root_disk` - (Optional, Bool, ForceNew) Whether to export only the system disk.
* `role_name` - (Optional, String, ForceNew) Role name (Default: CVM_QcsRole). Before exporting the images, make sure the role exists, and it has write permission to COS.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



