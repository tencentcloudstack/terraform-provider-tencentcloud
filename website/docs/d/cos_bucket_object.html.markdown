---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_bucket_object"
sidebar_current: "docs-tencentcloud-datasource-cos_bucket_object"
description: |-
  Use this data source to query the metadata of an object stored inside a bucket.
---

# tencentcloud_cos_bucket_object

Use this data source to query the metadata of an object stored inside a bucket.

## Example Usage

```hcl
data "tencentcloud_cos_bucket_object" "mycos" {
    bucket = "mycos-test-1258798060"
    key    = "hello-world.py"
    result_output_file  = "TFresults"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) Name of the bucket that contains the objects to query.
* `key` - (Required) The full path to the object inside the bucket.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cache_control` - Specifies caching behavior along the request/reply chain.
* `content_disposition` - Specifies presentational information for the object.
* `content_encoding` - Specifies what content encodings have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field.
* `content_type` - A standard MIME type describing the format of the object data.
* `etag` - ETag generated for the object，which is may not equal to MD5 value.
* `last_modified` - Last modified date of the object.
* `storage_class` - Object storage type such as STANDARD.


