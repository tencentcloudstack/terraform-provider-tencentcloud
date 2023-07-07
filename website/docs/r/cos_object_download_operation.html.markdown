---
subcategory: "Cloud Object Storage(COS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_object_download_operation"
sidebar_current: "docs-tencentcloud-resource-cos_object_download_operation"
description: |-
  Provides a resource to download object
---

# tencentcloud_cos_object_download_operation

Provides a resource to download object

## Example Usage

```hcl
resource "tencentcloud_cos_object_download_operation" "object_download" {
  bucket        = "xxxxxxx"
  key           = "test.txt"
  download_path = "/tmp/test.txt"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String, ForceNew) Bucket.
* `download_path` - (Required, String, ForceNew) Download path.
* `key` - (Required, String, ForceNew) Object key.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



