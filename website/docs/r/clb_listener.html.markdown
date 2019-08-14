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

```hcl
resource "tencentcloud_clb_listener" "tcp_listener" {
  clb_id                     = "lb-k2zjp9lv"
  listener_name              = "mylistener"
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
resource "tencentcloud_clb_listener" "https_listener" {
  clb_id               = "lb-k2zjp9lv"
  listener_name        = "listener_https"
  port                 = 80
  protocol             = "HTTPS"
  certificate_ssl_mode = "MUTUAL"
  certificate_id       = "mycert server ID "
  certificate_ca_id    = "mycert ca ID"
  sni_switch           = true
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required, ForceNew) ID of the CLB.
* `listener_name` - (Required) Name of the CLB listener, and available values can only be Chinese characters, English letters, numbers, underscore and hyphen '-'.
* `protocol` - (Required, ForceNew) Type of protocol within the listener, and available values include TCP, UDP, HTTP, HTTPS and TCP_SSL.
* `certificate_ca_id` - (Optional) ID of the client certificate. NOTES: Only supports listeners of 'HTTPS' and 'TCP_SSL' protocol.
* `certificate_id` - (Optional) ID of the server certificate. NOTES: Only supports listeners of 'HTTPS' and 'TCP_SSL' protocol.
* `certificate_ssl_mode` - (Optional) Type of certificate, and available values inclue 'UNIDIRECTIONAL', 'MUTUAL'. NOTES: Only supports listeners of 'HTTPS' and 'TCP_SSL' protocol.
* `health_check_health_num` - (Optional) Health threshold of health check, and the default is 3. If a success result is returned for the health check for 3 consecutive times, the backend CVM is identified as healthy. The value range is 2-10.
* `health_check_interval_time` - (Optional) Interval time of health check. The value range is 5-300 sec, and the default is 5 sec.
* `health_check_switch` - (Optional) Indicates whether health check is enabled.
* `health_check_time_out` - (Optional) Response timeout of health check. The value range is 2-60 sec, and the default is 2 sec. Response timeout needs to be less than check interval. NOTES: Only supports listeners of 'TCP','UDP','TCP_SSL' protocol.
* `health_check_unhealth_num` - (Optional) Unhealth threshold of health check, and the default is 3. If a success result is returned for the health check 3 consecutive times, the CVM is identified as unhealthy. The value range is 2-10.
* `port` - (Optional, ForceNew) Port of the CLB listener.
* `scheduler` - (Optional) Scheduling method of the CLB listener, and available values include 'WRR' and 'LEAST_CONN'. The default is 'WRR'. NOTES: The listener of HTTP and 'HTTPS' protocol additionally supports the 'IP Hash' method.
* `session_expire_time` - (Optional) Time of session persistence within the CLB listener. NOTES: Available when scheduler is specified as 'WRR', and not available when listener protocol is TCP_SSL.
* `sni_switch` - (Optional, ForceNew) Indicates whether SNI is enabled, and only supported with protocol 'HTTPS'.


