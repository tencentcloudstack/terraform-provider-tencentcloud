---
subcategory: "SSL Certificates(ssl)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_describe_host_cos_instance_list"
sidebar_current: "docs-tencentcloud-datasource-ssl_describe_host_cos_instance_list"
description: |-
  Use this data source to query detailed information of ssl describe_host_cos_instance_list
---

# tencentcloud_ssl_describe_host_cos_instance_list

Use this data source to query detailed information of ssl describe_host_cos_instance_list

## Example Usage

```hcl
data "tencentcloud_ssl_describe_host_cos_instance_list" "describe_host_cos_instance_list" {
  certificate_id = "8u8DII0l"
  resource_type  = "cos"
}
```

## Argument Reference

The following arguments are supported:

* `certificate_id` - (Required, String) Certificate ID to be deployed.
* `resource_type` - (Required, String) Deploy resource type cos.
* `filters` - (Optional, List) List of filter parameters.
* `is_cache` - (Optional, Int) Whether to query the cache, 1: Yes; 0: No, the default is the query cache, the cache is half an hour.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `filter_key` - (Required, String) Filter parameter key.
* `filter_value` - (Required, String) Filter parameter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `async_cache_time` - Current cache read timeNote: This field may return NULL, indicating that the valid value cannot be obtained.
* `async_offset` - Asynchronous refresh current execution numberNote: This field may return NULL, indicating that the valid value cannot be obtained.
* `async_total_num` - The total number of asynchronous refreshNote: This field may return NULL, indicating that the valid value cannot be obtained.
* `instance_list` - COS instance listNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `bucket` - Reserve bucket nameNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `cert_id` - Binded certificate IDNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `domain` - domain name.
  * `region` - Barrel areaNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `status` - Enabled: domain name online statusDisabled: Domain name offline status.


