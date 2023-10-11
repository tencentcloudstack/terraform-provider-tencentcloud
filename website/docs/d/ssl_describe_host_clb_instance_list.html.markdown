---
subcategory: "SSL Certificates"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_describe_host_clb_instance_list"
sidebar_current: "docs-tencentcloud-datasource-ssl_describe_host_clb_instance_list"
description: |-
  Use this data source to query detailed information of ssl describe_host_clb_instance_list
---

# tencentcloud_ssl_describe_host_clb_instance_list

Use this data source to query detailed information of ssl describe_host_clb_instance_list

## Example Usage

```hcl
data "tencentcloud_ssl_describe_host_clb_instance_list" "describe_host_clb_instance_list" {
  certificate_id = "8u8DII0l"
}
```

## Argument Reference

The following arguments are supported:

* `certificate_id` - (Required, String) Certificate ID to be deployed.
* `async_cache` - (Optional, Int) Whether to cache asynchronous.
* `filters` - (Optional, List) List of filtering parameters; Filterkey: domainmatch.
* `is_cache` - (Optional, Int) Whether to query the cache, 1: Yes; 0: No, the default is the query cache, the cache is half an hour.
* `old_certificate_id` - (Optional, String) Original certificate ID.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `filter_key` - (Required, String) Filter parameter key.
* `filter_value` - (Required, String) Filter parameter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `async_cache_time` - Current cache read timeNote: This field may return NULL, indicating that the valid value cannot be obtained.
* `async_offset` - Asynchronous refresh current execution numberNote: This field may return NULL, indicating that the valid value cannot be obtained.
* `async_total_num` - The total number of asynchronous refreshNote: This field may return NULL, indicating that the valid value cannot be obtained.
* `instance_list` - CLB instance listener listNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `listeners` - CLB listener listNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `certificate` - Certificate data binding of listenersNote: This field may return NULL, indicating that the valid value cannot be obtained.
      * `cert_ca_id` - Root certificate IDNote: This field may return NULL, indicating that the valid value cannot be obtained.
      * `cert_id` - Certificate ID.
      * `dns_names` - Domain name binding of certificates.
      * `s_s_l_mode` - Certificate certification mode: unidirectional unidirectional authentication, Mutual two -way certificationNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `listener_id` - Listener ID.
    * `listener_name` - Name of listeners.
    * `no_match_domains` - List of non -matching fieldsNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `protocol` - Type of listener protocol, https | TCP_SSL.
    * `rules` - List of listeners&#39; rulesNote: This field may return NULL, indicating that the valid value cannot be obtained.
      * `certificate` - Certificate data that has been bound to the rulesNote: This field may return NULL, indicating that the valid value cannot be obtained.
        * `cert_ca_id` - Root certificate IDNote: This field may return NULL, indicating that the valid value cannot be obtained.
        * `cert_id` - Certificate ID.
        * `dns_names` - Domain name binding of certificates.
        * `s_s_l_mode` - Certificate certification mode: unidirectional unidirectional authentication, Mutual two -way certificationNote: This field may return NULL, indicating that the valid value cannot be obtained.
      * `domain` - Domain name binding.
      * `is_match` - Whether the rules match the domain name to be bound to the certificate.
      * `location_id` - Rule ID.
      * `no_match_domains` - List of non -matching fieldsNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `sni_switch` - Whether to turn on SNI, 1 to open, 0 to close.
  * `load_balancer_id` - CLB instance ID.
  * `load_balancer_name` - CLB instance name name.


