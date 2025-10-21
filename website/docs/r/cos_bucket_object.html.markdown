---
subcategory: "Cloud Object Storage(COS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_bucket_object"
sidebar_current: "docs-tencentcloud-resource-cos_bucket_object"
description: |-
  Provides a COS object resource to put an object(content or file) to the bucket.
---

# tencentcloud_cos_bucket_object

Provides a COS object resource to put an object(content or file) to the bucket.

## Example Usage

### Uploading a file to a bucket

```hcl
resource "tencentcloud_cos_bucket_object" "myobject" {
  bucket = "mycos-1258798060"
  key    = "new_object_key"
  source = "path/to/file"
}
```

### Uploading a content to a bucket

```hcl
resource "tencentcloud_cos_bucket" "mycos" {
  bucket = "mycos-1258798060"
  acl    = "public-read"
}

resource "tencentcloud_cos_bucket_object" "myobject" {
  bucket  = tencentcloud_cos_bucket.mycos.bucket
  key     = "new_object_key"
  content = "the content that you want to upload."
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String, ForceNew) The name of a bucket to use. Bucket format should be [custom name]-[appid], for example `mycos-1258798060`.
* `key` - (Required, String, ForceNew) The name of the object once it is in the bucket.
* `acl` - (Optional, String) The canned ACL to apply. Available values include `private`, `public-read`, and `public-read-write`. Defaults to `private`.
* `cache_control` - (Optional, String) Specifies caching behavior along the request/reply chain. For further details, RFC2616 can be referred.
* `content_disposition` - (Optional, String) Specifies presentational information for the object.
* `content_encoding` - (Optional, String) Specifies what content encodings have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field.
* `content_type` - (Optional, String) A standard MIME type describing the format of the object data.
* `content` - (Optional, String) Literal string value to use as the object content, which will be uploaded as UTF-8-encoded text.
* `etag` - (Optional, String) The ETag generated for the object (an MD5 sum of the object content).
* `source` - (Optional, String) The path to the source file being uploaded to the bucket.
* `storage_class` - (Optional, String) Object storage type, Available values include `STANDARD_IA`, `MAZ_STANDARD_IA`, `INTELLIGENT_TIERING`, `MAZ_INTELLIGENT_TIERING`, `ARCHIVE`, `DEEP_ARCHIVE`. For more information, please refer to: https://cloud.tencent.com/document/product/436/33417.
* `tags` - (Optional, Map) Tag of the object.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



