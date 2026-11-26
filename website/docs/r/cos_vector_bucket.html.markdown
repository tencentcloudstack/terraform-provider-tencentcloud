---
subcategory: "Cloud Object Storage(COS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_vector_bucket"
sidebar_current: "docs-tencentcloud-resource-cos_vector_bucket"
description: |-
  Provides a resource to manage a COS vector bucket.
---

# tencentcloud_cos_vector_bucket

Provides a resource to manage a COS vector bucket.

## Example Usage

```hcl
data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

resource "tencentcloud_cos_vector_bucket" "example" {
  vector_bucket_name = "tf-example-${local.app_id}"

  encryption_configuration = {
    sse_type = "AES256"
  }
}
```

## Argument Reference

The following arguments are supported:

* `vector_bucket_name` - (Required, String) Name of the vector bucket, in the format of `<BucketName-APPID>`, e.g. `examplebucket-1250000000`. It supports lowercase letters, numbers and `-`, with a length of 3-63 characters. It cannot be changed after creation.
* `encryption_configuration` - (Optional, Object) Server-side encryption configuration of the vector bucket. Only SSE-COS is supported. It cannot be changed after creation.

The `encryption_configuration` object supports the following:

* `sse_type` - (Required, String) Server-side encryption type of the vector bucket. Only `AES256` (SSE-COS) is supported.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource. It is the same as `vector_bucket_name`.


## Import

COS vector bucket can be imported using the id, e.g.

```
terraform import tencentcloud_cos_vector_bucket.example tf-example-1309101522
```

