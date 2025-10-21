---
subcategory: "SSL Certificates(ssl)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_describe_host_tke_instance_list"
sidebar_current: "docs-tencentcloud-datasource-ssl_describe_host_tke_instance_list"
description: |-
  Use this data source to query detailed information of ssl describe_host_tke_instance_list
---

# tencentcloud_ssl_describe_host_tke_instance_list

Use this data source to query detailed information of ssl describe_host_tke_instance_list

## Example Usage

```hcl
data "tencentcloud_ssl_describe_host_tke_instance_list" "describe_host_tke_instance_list" {
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
  * `cluster_id` - Cluster ID.
  * `cluster_name` - Cluster name.
  * `cluster_type` - Cluster.
  * `cluster_version` - Cluster.
  * `namespace_list` - Cluster Naming Space List.
    * `name` - namespace name.
    * `secret_list` - Secret list.
      * `cert_id` - Certificate ID.
      * `ingress_list` - Ingress list.
        * `domains` - Ingress domain name list.
        * `ingress_name` - Ingress name.
        * `tls_domains` - TLS domain name list.
      * `name` - Secret name.
      * `no_match_domains` - List of domain names that are not matched with the new certificateNote: This field may return NULL, indicating that the valid value cannot be obtained.


