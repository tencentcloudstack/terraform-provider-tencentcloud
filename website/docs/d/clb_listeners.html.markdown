---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_listeners"
sidebar_current: "docs-tencentcloud-datasource-clb_listeners"
description: |-
  Use this data source to query detailed information of CLB listener
---

# tencentcloud_clb_listeners

Use this data source to query detailed information of CLB listener

## Example Usage

```hcl
data "tencentcloud_clb_listeners" "foo" {
  clb_id      = "lb-k2zjp9lv"
  listener_id = "lbl-mwr6vbtv"
  protocol    = "TCP"
  port        = 80
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required) Id of the CLB to be queried.
* `listener_id` - (Optional) Id of the listener to be queried.
* `port` - (Optional) Port of the CLB listener.
* `protocol` - (Optional) Type of protocol within the listener, and available values are 'TCP', 'UDP', 'HTTP', 'HTTPS' and 'TCP_SSL'.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `listener_list` - A list of listeners of cloud load balancers. Each element contains the following attributes:
  * `certificate_ca_id` - Id of the client certificate. It must be set when SSLMode is 'mutual'. NOTES: only supported by listeners of 'HTTPS' and 'TCP_SSL' protocol.
  * `certificate_id` - Id of the server certificate. It must be set when protocol is 'HTTPS' or 'TCP_SSL'. NOTES: only supported by listeners of 'HTTPS' and 'TCP_SSL' protocol and must be set when it is available.
  * `certificate_ssl_mode` - Type of certificate, and available values inclue 'UNIDIRECTIONAL', 'MUTUAL'. NOTES: Only supports listeners of 'HTTPS' and 'TCP_SSL' protocol and must be set when it is available.
  * `clb_id` - Id of the CLB.
  * `health_check_health_num` - Health threshold of health check, and the default is 3. If a success result is returned for the health check three consecutive times, the CVM is identified as healthy. The value range is 2-10. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.
  * `health_check_interval_time` - Interval time of health check. The value range is 5-300 sec, and the default is 5 sec. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.
  * `health_check_switch` - Indicates whether health check is enabled.
  * `health_check_time_out` - Response timeout of health check. The value range is 2-60 sec, and the default is 2 sec. Response timeout needs to be less than check interval. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration.
  * `health_check_unhealth_num` - Unhealthy threshold of health check, and the default is 3. If a success result is returned for the health check three consecutive times, the CVM is identified as unhealthy. The value range is 2-10. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.
  * `listener_id` - Id of the listener.
  * `listener_name` - Name of the CLB listener.
  * `port` - Port of the CLB listener.
  * `protocol` - Protocol of the listener. Available values are 'HTTP', 'HTTPS', 'TCP', 'UDP', 'TCP_SSL'.
  * `scheduler` - Scheduling method of the CLB listener, and available values are 'WRR' and 'LEAST_CONN'. The default is 'WRR'. NOTES: The listener of 'HTTP' and 'HTTPS' protocol additionally supports the 'IP HASH' method. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.
  * `session_expire_time` - Time of session persistence within the CLB listener. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.
  * `sni_switch` - Indicates whether SNI is enabled. NOTES: Only supported by 'HTTPS' protocol.


