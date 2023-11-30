Provides a COS resource to create a COS bucket policy and set its attributes.

Example Usage

```hcl
resource "tencentcloud_cos_bucket_policy" "cos_policy" {
  bucket = "mycos-1258798060"

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
        "qcs::cos:<bucket region>:uid/<your-account-id>:<bucket name>/*"
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
$ terraform import tencentcloud_cos_bucket_policy.bucket bucket-name
```