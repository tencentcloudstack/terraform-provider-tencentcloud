---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_listener"
sidebar_current: "docs-tencentcloud-resource-clb_listener"
description: |-
  Provides a resource to create a CLB listener.
---

# tencentcloud_clb_listener

Provides a resource to create a CLB listener.

## Example Usage

HTTP Listener

```hcl
resource "tencentcloud_clb_listener" "HTTP_listener" {
  clb_id        = "lb-0lh5au7v"
  listener_name = "test_listener"
  port          = 80
  protocol      = "HTTP"
}
```

TCP/UDP Listener

```hcl
resource "tencentcloud_clb_listener" "TCP_listener" {
  clb_id                     = "lb-0lh5au7v"
  listener_name              = "test_listener"
  port                       = 80
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 2
  health_check_interval_time = 5
  health_check_health_num    = 3
  health_check_unhealth_num  = 3
  session_expire_time        = 30
  scheduler                  = "WRR"
}
```

HTTPS Listener

```hcl
resource "tencentcloud_clb_listener" "HTTPS_listener" {
  clb_id               = "lb-0lh5au7v"
  listener_name        = "test_listener"
  port                 = "80"
  protocol             = "HTTPS"
  certificate_ssl_mode = "MUTUAL"
  certificate_id       = "VjAYq9xc"
  certificate_ca_id    = "VfqcL1ME"
  sni_switch           = true
}
```

TCP SSL Listener

```hcl
resource "tencentcloud_clb_listener" "TCPSSL_listener" {
  clb_id                     = "lb-0lh5au7v"
  listener_name              = "test_listener"
  port                       = "80"
  protocol                   = "TCP_SSL"
  certificate_ssl_mode       = "MUTUAL"
  certificate_id             = "VjAYq9xc"
  certificate_ca_id          = "VfqcL1ME"
  health_check_switch        = true
  health_check_time_out      = 2
  health_check_interval_time = 5
  health_check_health_num    = 3
  health_check_unhealth_num  = 3
  scheduler                  = "WRR"
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required, ForceNew) Id of the CLB.
* `listener_name` - (Required) Name of the CLB listener, and available values can only be Chinese characters, English letters, numbers, underscore and hyphen '-'.
* `protocol` - (Required, ForceNew) Type of protocol within the listener, and available values include 'TCP', 'UDP', 'HTTP', 'HTTPS' and 'TCP_SSL'.
* `certificate_ca_id` - (Optional) Id of the client certificate. NOTES: Only supports listeners of 'HTTPS' and 'TCP_SSL' protocol and must be set when the ssl mode is 'MUTUAL'.
* `certificate_id` - (Optional) Id of the server certificate. NOTES: Only supports listeners of 'HTTPS' and 'TCP_SSL' protocol and must be set when it is available.
* `certificate_ssl_mode` - (Optional) Type of certificate, and available values inclue 'UNIDIRECTIONAL', 'MUTUAL'. NOTES: Only supports listeners of 'HTTPS' and 'TCP_SSL' protocol and must be set when it is available.
* `health_check_health_num` - (Optional) Health threshold of health check, and the default is 3. If a success result is returned for the health check for 3 consecutive times, the backend CVM is identified as healthy. The value range is 2-10. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.
* `health_check_interval_time` - (Optional) Interval time of health check. The value range is 5-300 sec, and the default is 5 sec. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.
* `health_check_switch` - (Optional) Indicates whether health check is enabled.
* `health_check_time_out` - (Optional) Response timeout of health check. The value range is 2-60 sec, and the default is 2 sec. Response timeout needs to be less than check interval. NOTES: Only supports listeners of 'TCP','UDP','TCP_SSL' protocol.
* `health_check_unhealth_num` - (Optional) Unhealth threshold of health check, and the default is 3. If a success result is returned for the health check 3 consecutive times, the CVM is identified as unhealthy. The value range is 2-10. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.
* `port` - (Optional, ForceNew) Port of the CLB listener.
* `scheduler` - (Optional) Scheduling method of the CLB listener, and available values include 'WRR' and 'LEAST_CONN'. The default is 'WRR'. NOTES: The listener of HTTP and 'HTTPS' protocol additionally supports the 'IP Hash' method. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.
* `session_expire_time` - (Optional) Time of session persistence within the CLB listener. NOTES: Available when scheduler is specified as 'WRR', and not available when listener protocol is 'TCP_SSL'. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.
* `sni_switch` - (Optional, ForceNew) Indicates whether SNI is enabled, and only supported with protocol 'HTTPS'. If enabled, you can set a certificate for each rule in `tencentcloud_clb_listener_rule`, otherwise all rules have a certificate.


