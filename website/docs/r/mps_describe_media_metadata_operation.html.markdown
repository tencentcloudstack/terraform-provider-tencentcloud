---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_describe_media_metadata_operation"
sidebar_current: "docs-tencentcloud-resource-mps_describe_media_metadata_operation"
description: |-
  Provides a resource to create a mps describe_media_metadata_operation
---

# tencentcloud_mps_describe_media_metadata_operation

Provides a resource to create a mps describe_media_metadata_operation

## Example Usage

### Operation through COS

```hcl
data "tencentcloud_cos_bucket_object" "object" {
  bucket = "keep-bucket-${local.app_id}"
  key    = "/mps-test/test.mov"
}

resource "tencentcloud_mps_describe_media_metadata_operation" "operation" {
  input_info {
    type = "COS"
    cos_input_info {
      bucket = data.tencentcloud_cos_bucket_object.object.bucket
      region = "%s"
      object = data.tencentcloud_cos_bucket_object.object.key
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `input_info` - (Required, List, ForceNew) Input information of file for metadata getting.

The `cos_input_info` object supports the following:

* `bucket` - (Required, String) The COS bucket of the object to process, such as `TopRankVideo-125xxx88`.
* `object` - (Required, String) The path of the object to process, such as `/movie/201907/WildAnimal.mov`.
* `region` - (Required, String) The region of the COS bucket, such as `ap-chongqing`.

The `input_info` object supports the following:

* `type` - (Required, String) The input type. Valid values: `COS`: A COS bucket address.  `URL`: A URL. `AWS-S3`: An AWS S3 bucket address. Currently, this type is only supported for transcoding tasks.
* `cos_input_info` - (Optional, List) The information of the COS object to process. This parameter is valid and required when `Type` is `COS`.
* `s3_input_info` - (Optional, List) The information of the AWS S3 object processed. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.
* `url_input_info` - (Optional, List) The URL of the object to process. This parameter is valid and required when `Type` is `URL`.Note: This field may return null, indicating that no valid value can be obtained.

The `s3_input_info` object supports the following:

* `s3_bucket` - (Required, String) The AWS S3 bucket.
* `s3_object` - (Required, String) The path of the AWS S3 object.
* `s3_region` - (Required, String) The region of the AWS S3 bucket.
* `s3_secret_id` - (Optional, String) The key ID required to access the AWS S3 object.
* `s3_secret_key` - (Optional, String) The key required to access the AWS S3 object.

The `url_input_info` object supports the following:

* `url` - (Required, String) URL of a video.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



