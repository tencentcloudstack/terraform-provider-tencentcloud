---
subcategory: "SSL Certificates"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_describe_host_waf_instance_list"
sidebar_current: "docs-tencentcloud-datasource-ssl_describe_host_waf_instance_list"
description: |-
  Use this data source to query detailed information of ssl describe_host_waf_instance_list
---

# tencentcloud_ssl_describe_host_waf_instance_list

Use this data source to query detailed information of ssl describe_host_waf_instance_list

## Example Usage

```hcl
data "tencentcloud_ssl_describe_host_waf_instance_list" "describe_host_waf_instance_list" {
  certificate_id = "8u8DII0l"
  resource_type  = "waf"
}
```

## Argument Reference

The following arguments are supported:

* `certificate_id` - (Required, String) Certificate ID to be deployed.
* `resource_type` - (Required, String) Deploy resource type.
* `filters` - (Optional, List) List of filtering parameters; Filterkey: domainmatch.
* `is_cache` - (Optional, Int) Whether to query the cache, 1: Yes; 0: No, the default is the query cache, the cache is half an hour.
* `old_certificate_id` - (Optional, String) Deployed certificate ID.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `filter_key` - (Required, String) Filter parameter key.
* `filter_value` - (Required, String) Filter parameter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - WAF instance listNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `cert_id` - Binded certificate IDNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `domain` - domain name.
  * `status` - -1: Unrelated certificate of domain name.1: The domain name HTTPS has been opened.0: The domain name HTTPS has been closed.


