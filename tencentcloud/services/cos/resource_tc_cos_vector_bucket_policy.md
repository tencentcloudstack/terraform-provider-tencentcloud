Provides a resource to manage a COS vector bucket policy.

~> **NOTE:** All policies must be configured within a single `tencentcloud_cos_vector_bucket_policy` resource instance. Creating multiple instances of this resource to manage individual policies is not supported.

Example Usage

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

Import

COS vector bucket policy can be imported using the id, e.g.

```
terraform import tencentcloud_cos_vector_bucket_policy.example examplebucket-1309101522
```
