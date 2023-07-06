---
subcategory: "Cloud Object Storage(COS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_bucket_multipart_uploads"
sidebar_current: "docs-tencentcloud-datasource-cos_bucket_multipart_uploads"
description: |-
  Use this data source to query the COS bucket multipart uploads.
---

# tencentcloud_cos_bucket_multipart_uploads

Use this data source to query the COS bucket multipart uploads.

## Example Usage

```hcl
data "tencentcloud_cos_bucket_multipart_uploads" "cos_bucket_multipart_uploads" {
  bucket = "xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String) Bucket.
* `delimiter` - (Optional, String) The delimiter is a symbol, and the Object name contains the Object between the specified prefix and the first occurrence of delimiter characters as a set of elements: common prefix. If there is no prefix, start from the beginning of the path.
* `encoding_type` - (Optional, String) Specifies the encoding format of the return value. Legal value: url.
* `prefix` - (Optional, String) The returned Object key must be prefixed with Prefix. Note that when using the prefix query, the returned key still contains Prefix.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `uploads` - Information for each Upload.
  * `initiated` - The starting time of multipart upload.
  * `initiator` - Used to represent the information of the initiator of this upload.
    * `display_name` - Abbreviation for user identity ID (UIN).
    * `id` - The user's unique CAM identity ID.
  * `key` - Name of the Object.
  * `owner` - Information used to represent the owner of these chunks.
    * `display_name` - Abbreviation for user identity ID (UIN).
    * `id` - The user's unique CAM identity ID.
  * `storage_class` - Used to represent the storage level of a chunk. Enumerated value: STANDARD,STANDARD_IA,ARCHIVE.
  * `upload_id` - Mark the ID of this multipart upload.


