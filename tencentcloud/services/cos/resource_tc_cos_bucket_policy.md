Provides a COS resource to create a COS bucket policy and set its attributes.

Example Usage

```hcl
data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

resource "tencentcloud_cos_bucket" "example" {
  bucket = "private-bucket-${local.app_id}"
  acl    = "private"
}

resource "tencentcloud_cos_bucket_policy" "example" {
  bucket = tencentcloud_cos_bucket.example.id
  policy = <<EOF
{
  "version": "2.0",
  "Statement": [
    {
      "Principal": {
        "qcs": [
          "qcs::cam::uin/<your-account-id>:uin/<your-account-id>"
        ]
      },
      "Action": [
        "name/cos:DeleteBucket",
        "name/cos:PutBucketACL"
      ],
      "Effect": "allow",
      "Resource": [
        "qcs::cos:<bucket region>:uid/<your-appid-id>:<your-bucket-name>/*"
      ]
    }
  ]
}
EOF
}
```

Import

COS bucket policy can be imported, e.g.

```
$ terraform import tencentcloud_cos_bucket_policy.example private-bucket-1309118521
```