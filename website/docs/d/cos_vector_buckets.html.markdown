---
subcategory: "Cloud Object Storage(COS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_vector_buckets"
sidebar_current: "docs-tencentcloud-data_source-cos_vector_buckets"
description: |-
  Use this data source to query the list of COS vector buckets.
---

# tencentcloud_cos_vector_buckets

Use this data source to query the list of COS vector buckets.

## Example Usage

### Query all buckets

```hcl
data "tencentcloud_cos_vector_buckets" "example" {}
```

### Query buckets by prefix

```hcl
data "tencentcloud_cos_vector_buckets" "example" {
  prefix = "tf-example"
}
```

## Argument Reference

The following arguments are supported:

* `prefix` - (Optional, String) Limits the response to vector buckets whose names begin with the specified prefix, 1-63 characters in length.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `creation_time` - Creation date and time of the vector bucket, a Unix timestamp.
  * `vector_bucket_name` - Name of the vector bucket, in the format of `<BucketName-APPID>`, e.g. `examplebucket-1250000000`.
  * `vector_bucket_qcs` - QCS (resource name) of the vector bucket, e.g. `qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/examplebucket-1250000000`.
* `id` - ID of the data source, a random token.
* `vector_buckets` - List of vector buckets.


