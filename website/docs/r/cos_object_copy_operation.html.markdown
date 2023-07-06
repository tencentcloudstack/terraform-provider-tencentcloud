---
subcategory: "Cloud Object Storage(COS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_object_copy_operation"
sidebar_current: "docs-tencentcloud-resource-cos_object_copy_operation"
description: |-
  Provides a resource to copy object
---

# tencentcloud_cos_object_copy_operation

Provides a resource to copy object

## Example Usage

```hcl
resource "tencentcloud_cos_object_copy_operation" "object_copy" {
  bucket     = "keep-copy-xxxxxxx"
  key        = "copy-acl.txt"
  source_url = "keep-test-xxxxxx.cos.ap-guangzhou.myqcloud.com/acl.txt"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String, ForceNew) Bucket.
* `key` - (Required, String, ForceNew) Object key.
* `source_url` - (Required, String, ForceNew) Object key.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



