---
subcategory: "Cloud Object Storage(COS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_vector_indexes"
sidebar_current: "docs-tencentcloud-data_source-cos_vector_indexes"
description: |-
  Use this data source to query the list of indexes of a COS vector bucket.
---

# tencentcloud_cos_vector_indexes

Use this data source to query the list of indexes of a COS vector bucket.

## Example Usage

```hcl
data "tencentcloud_cos_vector_indexes" "example" {
  vector_bucket_name = "tf-example-1309101522"
  prefix             = "tf-example"
}
```

## Argument Reference

The following arguments are supported:

* `vector_bucket_name` - (Required, String) Name of the vector bucket whose indexes are queried, in the format of `<BucketName-APPID>`, e.g. `examplebucket-1250000000`.
* `prefix` - (Optional, String) Limits the response to indexes whose names begin with the specified prefix, 1-63 characters in length.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `creation_time` - Creation date and time of the index, a Unix timestamp.
  * `index_name` - Name of the index, 3-63 characters in length.
  * `index_qcs` - QCS (resource name) of the index, e.g. `qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/examplebucket-1250000000/index/exampleindex`.
  * `vector_bucket_name` - Name of the vector bucket the index belongs to, in the format of `<BucketName-APPID>`, e.g. `examplebucket-1250000000`.
* `id` - ID of the data source, a random token.
* `indexes` - List of indexes.


