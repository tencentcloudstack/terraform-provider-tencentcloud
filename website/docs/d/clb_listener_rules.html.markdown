---
subcategory: "Cloud Load Balancer(CLB)"
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
  rule_id     = "loc-inem40hz"
  domain      = "abc.com"
  url         = "/"
  scheduler   = "WRR"
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required, String) ID of the CLB to be queried.
* `listener_id` - (Required, String) ID of the CLB listener to be queried.
* `domain` - (Optional, String) Domain name of the forwarding rule to be queried.
* `result_output_file` - (Optional, String) Used to save results.
* `rule_id` - (Optional, String) ID of the forwarding rule to be queried.
* `scheduler` - (Optional, String) Scheduling method of the forwarding rule of thr CLB listener, and available values include `WRR`, `IP HASH` and `LEAST_CONN`. The default is `WRR`.
* `url` - (Optional, String) Url of the forwarding rule to be queried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `rule_list` - A list of forward rules of listeners. Each element contains the following attributes:
  * `certificate_ca_id` - ID of the client certificate. NOTES: Only supports listeners of 'HTTPS' and 'TCP_SSL' protocol.
  * `certificate_id` - ID of the server certificate. NOTES: Only supports listeners of 'HTTPS'  and 'TCP_SSL' protocol.
  * `certificate_ssl_mode` - Type of SSL Mode, and available values inclue 'UNIDIRECTIONAL', 'MUTUAL'.NOTES: Only supports listeners of 'HTTPS'  and 'TCP_SSL' protocol.
  * `clb_id` - ID of the CLB.
  * `health_check_health_num` - Health threshold of health check, and the default is `3`. If a success result is returned for the health check three consecutive times, the CVM is identified as healthy. The value range is 2-10. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.
  * `health_check_http_code` - HTTP Status Code. The default is 31 and value range is 1-31. 1 means the return value '1xx' is health. 2 means the return value '2xx' is health. 4 means the return value '3xx' is health. 8 means the return value 4xx is health. 16 means the return value '5xx' is health. If you want multiple return codes to indicate health, need to add the corresponding values. NOTES: The 'HTTP' health check of the 'TCP' listener only supports specifying one health check status code. NOTES: Only supports listeners of 'HTTP' and 'HTTPS' protocol.
  * `health_check_http_domain` - Domain name of health check. NOTES: Only supports listeners of 'HTTPS' and 'HTTP' protocol.
  * `health_check_http_method` - Methods of health check. NOTES: Only supports listeners of 'HTTPS' and 'HTTP' protocol. The default is 'HEAD', the available value include 'HEAD' and 'GET'.
  * `health_check_http_path` - Path of health check. NOTES: Only supports listeners of 'HTTPS' and 'HTTP' protocol.
  * `health_check_interval_time` - Interval time of health check. The value range is 2-300 sec, and the default is `5` sec. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.
  * `health_check_switch` - Indicates whether health check is enabled.
  * `health_check_unhealth_num` - Unhealth threshold of health check, and the default is `3`. If a success result is returned for the health check three consecutive times, the CVM is identified as unhealthy. The value range is 2-10. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.
  * `http2_switch` - Indicate to set HTTP2 protocol or not.
  * `listener_id` - ID of the listener.
  * `rule_id` - ID of the rule.
  * `scheduler` - Scheduling method of the CLB listener, and available values include 'WRR', 'IP_HASH' and 'LEAST_CONN'. The default is 'WRR'. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.
  * `session_expire_time` - Time of session persistence within the CLB listener. NOTES: Available when scheduler is specified as 'WRR'. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.


