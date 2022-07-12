---
subcategory: "Cloud Object Storage(COS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_buckets"
sidebar_current: "docs-tencentcloud-datasource-cos_buckets"
description: |-
  Use this data source to query the COS buckets of the current Tencent Cloud user.
---

# tencentcloud_cos_buckets

Use this data source to query the COS buckets of the current Tencent Cloud user.

## Example Usage

```hcl
data "tencentcloud_cos_buckets" "cos_buckets" {
  bucket_prefix      = "tf-bucket-"
  result_output_file = "mytestpath"
}
```

## Argument Reference

The following arguments are supported:

* `bucket_prefix` - (Optional, String) A prefix string to filter results by bucket name.
* `result_output_file` - (Optional, String) Used to save results.
* `tags` - (Optional, Map) Tags to filter bucket.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `bucket_list` - A list of bucket. Each element contains the following attributes:
  * `acl_body` - Bucket verbose acl configurations.
  * `acl` - Bucket access control configurations.
  * `bucket` - Bucket name, the format likes `<bucket>-<appid>`.
  * `cors_rules` - A list of CORS rule configurations.
    * `allowed_headers` - Specifies which headers are allowed.
    * `allowed_methods` - Specifies which methods are allowed. Can be GET, PUT, POST, DELETE or HEAD.
    * `allowed_origins` - Specifies which origins are allowed.
    * `expose_headers` - Specifies expose header in the response.
    * `max_age_seconds` - Specifies time in seconds that browser can cache the response for a preflight request.
  * `cos_bucket_url` - The URL of this cos bucket.
  * `lifecycle_rules` - The lifecycle configuration of a bucket.
    * `expiration` - Specifies a period in the object's expire.
      * `date` - Specifies the date after which you want the corresponding action to take effect.
      * `days` - Specifies the number of days after object creation when the specific rule action takes effect.
    * `filter_prefix` - Object key prefix identifying one or more objects to which the rule applies.
    * `non_current_expiration` - Specifies when non current object versions shall expire.
      * `non_current_days` - Number of days after non current object creation when the specific rule action takes effect. The maximum value is 3650.
    * `non_current_transition` - Specifies when to transition objects of non current versions and the target storage class.
      * `non_current_days` - Number of days after non current object creation when the specific rule action takes effect.
      * `storage_class` - Specifies the storage class to which you want the non current object to transition. Available values include STANDARD, STANDARD_IA and ARCHIVE.
    * `transition` - Specifies a period in the object's transitions.
      * `date` - Specifies the date after which you want the corresponding action to take effect.
      * `days` - Specifies the number of days after object creation when the specific rule action takes effect.
      * `storage_class` - Specifies the storage class to which you want the object to transition. Available values include STANDARD, STANDARD_IA and ARCHIVE.
  * `origin_domain_rules` - Bucket origin domain rules.
  * `origin_pull_rules` - Bucket Origin-Pull rules.
  * `tags` - The tags of a bucket.
  * `website` - A list of one element containing configuration parameters used when the bucket is used as a website.
    * `error_document` - An absolute path to the document to return in case of a 4XX error.
    * `index_document` - COS returns this index document when requests are made to the root domain or any of the subfolders.


