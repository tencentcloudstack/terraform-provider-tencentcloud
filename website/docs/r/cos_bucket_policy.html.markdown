---
subcategory: "Cloud Object Storage(COS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_bucket_policy"
sidebar_current: "docs-tencentcloud-resource-cos_bucket_policy"
description: |-
  Provides a COS resource to create a COS bucket policy and set its attributes.
---

# tencentcloud_cos_bucket_policy

Provides a COS resource to create a COS bucket policy and set its attributes.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String, ForceNew) The name of a bucket to be created. Bucket format should be [custom name]-[appid], for example `mycos-1258798060`.
* `policy` - (Required, String) The text of the policy. For more info please refer to [Tencent official doc](https://intl.cloud.tencent.com/document/product/436/18023), The six-segment resource scenario example in the document is as follows: Example of specifying a bucket: `qcs::cos:ap-guangzhou:uid/1250000000:examplebucket-1250000000/*`, In the cdc scenario: `qcs::cos:ap-guangzhou:uid/1250000000:cdc_cluster-123456_examplebucket-1250000000/*`; Example of specifying a folder: `qcs::cos:ap-guangzhou:uid/1250000000:examplebucket-1250000000/folder/*`, In the cdc scenario: `qcs::cos:ap-guangzhou:uid/1250000000:cdc_cluster-123456_examplebucket-1250000000/folder/*`; Specified object example: `qcs::cos:ap-guangzhou:uid/1250000000:examplebucket-1250000000/folder/exampleobject`, In the cdc scenario: `qcs::cos:ap-guangzhou:uid/1250000000:cdc_cluster-123456_examplebucket-1250000000/folder/exampleobject`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

COS bucket policy can be imported, e.g.

```
$ terraform import tencentcloud_cos_bucket_policy.example private-bucket-1309118521
```

