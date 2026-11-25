---
subcategory: "Cloud Object Storage(COS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_vector_bucket_policy"
sidebar_current: "docs-tencentcloud-resource-cos_vector_bucket_policy"
description: |-
  Provides a resource to manage a COS vector bucket policy.
---

# tencentcloud_cos_vector_bucket_policy

Provides a resource to manage a COS vector bucket policy.

~> **NOTE:** All policies must be configured within a single `tencentcloud_cos_vector_bucket_policy` resource instance. Creating multiple instances of this resource to manage individual policies is not supported.

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

resource "tencentcloud_cos_vector_bucket_policy" "example" {
  vector_bucket_name = tencentcloud_cos_vector_bucket.example.id
  policy = jsonencode({
    Version = "2.0"
    Statement = [
      {
        Sid    = "tf-example1"
        Effect = "Allow"
        Action = [
          "name/cos:DeleteVectorBucket",
          "name/cos:GetVectorBucket",
          "name/cos:ListVectorIndexes",
        ]
        Resource = [
          "qcs::cosvector:ap-guangzhou:uid/1309118522:bucket/tf-example-1309118522/*",
        ]
        Principal = {
          qcs = [
            "qcs::cam::uin/100023201586:uin/100023201586",
          ]
        }
      },
      {
        Sid    = "tf-example2"
        Effect = "Deny"
        Action = [
          "name/cos:PutVectorIndex",
          "name/cos:ListVectorIndexes",
        ]
        Resource = [
          "qcs::cosvector:ap-guangzhou:uid/1309118522:bucket/tf-example-1309118522/*",
        ]
        Principal = {
          qcs = [
            "qcs::cam::uin/100019282779:uin/100019282779",
          ]
        }
      }
    ]
  })
}
```

## Argument Reference

The following arguments are supported:

* `policy` - (Required, String) Policy of the vector bucket. It is a JSON-formatted access control policy serialized as a string. The new policy fully overwrites the old one.
* `vector_bucket_name` - (Required, String) Name of the vector bucket the policy attaches to, in the format of `<BucketName-APPID>`, e.g. `examplebucket-1250000000`. It cannot be changed after creation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource. It is the same as `vector_bucket_name`.


## Import

COS vector bucket policy can be imported using the id, e.g.

```
terraform import tencentcloud_cos_vector_bucket_policy.example examplebucket-1309101522
```

