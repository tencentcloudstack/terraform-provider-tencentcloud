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

### Query all cos buckets

```hcl
data "tencentcloud_cos_buckets" "example" {}
```

### Query cos buckets by filters

```hcl
data "tencentcloud_cos_buckets" "example" {
  bucket_prefix = "tf-example-prefix"
  tags = {
    createBy = "Terraform"
  }
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
    * `abort_incomplete_multipart_upload` - Set the maximum time a multipart upload is allowed to remain running.
      * `days_after_initiation` - Specifies the number of days after the multipart upload starts that the upload must be completed. The maximum value is 3650.
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
    * `domain` - Specify domain host.
    * `status` - Domain status, default: `ENABLED`.
    * `type` - Specify origin domain type, available values: `REST`, `WEBSITE`, `ACCELERATE`, default: `REST`.
  * `origin_pull_rules` - Bucket Origin-Pull rules.
    * `back_to_source_mode` - Back to source mode. Allow value: Proxy, Mirror, Redirect.
    * `custom_http_headers` - Specifies the custom headers that you can add for COS to access your origin server.
    * `follow_http_headers` - Specifies the pass through headers when accessing the origin server.
    * `follow_query_string` - Specifies whether to pass through COS request query string when accessing the origin server.
    * `follow_redirection` - Specifies whether to follow 3XX redirect to another origin server to pull data from.
    * `host` - Allows only a domain name or IP address. You can optionally append a port number to the address.
    * `prefix` - Triggers the origin-pull rule when the requested file name matches this prefix.
    * `priority` - Priority of origin-pull rules, do not set the same value for multiple rules.
    * `protocol` - the protocol used for COS to access the specified origin server. The available value include `HTTP`, `HTTPS` and `FOLLOW`.
    * `sync_back_to_source` - If `true`, COS will not return 3XX status code when pulling data from an origin server. Currently available zone: ap-beijing, ap-shanghai, ap-singapore, ap-mumbai.
  * `tags` - The tags of a bucket.
  * `website` - A list of one element containing configuration parameters used when the bucket is used as a website.
    * `error_document` - An absolute path to the document to return in case of a 4XX error.
    * `index_document` - COS returns this index document when requests are made to the root domain or any of the subfolders.


