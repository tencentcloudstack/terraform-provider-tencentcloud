---
subcategory: "Cloud Object Storage(COS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_object_abort_multipart_upload_operation"
sidebar_current: "docs-tencentcloud-resource-cos_object_abort_multipart_upload_operation"
description: |-
  Provides a resource to abort multipart upload
---

# tencentcloud_cos_object_abort_multipart_upload_operation

Provides a resource to abort multipart upload

## Example Usage

```hcl
resource "tencentcloud_cos_object_abort_multipart_upload_operation" "abort_multipart_upload" {
  bucket    = "keep-test-xxxxxx"
  key       = "object"
  upload_id = "xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String, ForceNew) Bucket.
* `key` - (Required, String, ForceNew) Object key.
* `upload_id` - (Required, String, ForceNew) Multipart uploaded id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



