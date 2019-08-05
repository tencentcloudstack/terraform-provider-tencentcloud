---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_listener_rules"
sidebar_current: "docs-tencentcloud-datasource-clb_listener_rules"
description: |-
  Use this data source to query detailed information of CLB listener rule
---

# tencentcloud_clb_listener_rules

Use this data source to query detailed information of CLB listener rule

## Example Usage

```hcl
data "tencentcloud_clb_listener_rules" "foo" {
  clb_id      = "lb-k2zjp9lv"
  listener_id = "lbl-mwr6vbtv"
  location_id = "loc-inem40hz"
  domain      = "abc.com"
  url         = "/"
  scheduler   = "WRR"
}
```

## Argument Reference

The following arguments are supported:

* `listener_id` - (Required) ID of the listener to be queried.
* `clb_id` - (Optional) ID of the CLB.
* `domain` - (Optional) Domain of the rule.
* `location_id` - (Optional) ID of the rule.
* `result_output_file` - (Optional) Used to save results.
* `scheduler` - (Optional) Scheduling method of the CLB listener, and available values include 'WRR' and 'LEAST_CONN'. The defaule is 'WRR'.
* `url` - (Optional) Url of the rule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `rule_list` - A list of forward rules of listeners. Each element contains the following attributes:
  * `certificate_ca_id` - ID of the client certificate. If not specified, the content, key, name of client certificate must be set when SSLMode is 'mutual'. NOTES: only supported by listeners of protocol 'HTTPS'.
  * `certificate_id` - ID of the server certificate. If not specified, the content, key, and name of the server certificate must be set. NOTES: only supported by listeners of protocol 'HTTPS'.
  * `certificate_ssl_mode` - Type of SSL Mode, and available values inclue 'UNIDRECTIONAL', 'MUTUAL'.
  * `clb_id` - ID of the CLB.
  * `health_check_health_num` - Health threshold of health check, and the default is 3. If a success result is returned for the health check three consecutive times, the CVM is identified as healthy. The value range is 2-10.
  * `health_check_http_code` - Path of health check (applicable only to HTTP/HTTPS check methods).
  * `health_check_http_domain` - Domain name of health check (applicable only to HTTP/HTTPS check methods)
  * `health_check_http_method` - Methods of health check (applicable only to HTTP/HTTPS check methods). Available values include 'HEAD' and 'GET'.
  * `health_check_http_path` - Path of health check (applicable only to HTTP/HTTPS check methods). 
  * `health_check_interval_time` - Interval time of health check. The value range is 5-300 sec, and the default is 5 sec.
  * `health_check_switch` - Indicates whether health check is enabled.
  * `health_check_unhealth_num` - Unhealth threshold of health check, and the default is 3. If a success result is returned for the health check three consecutive times, the CVM is identified as unhealthy. The value range is 2-10.
  * `listener_id` - ID of the listener.
  * `location_id` - ID of the rule.
  * `scheduler` - Scheduling method of the CLB listener, and available values include 'WRR' and 'LEAST_CONN'. The defaule is 'WRR'.
  * `session_expire_time` - Time of session persistence within the CLB listener.


