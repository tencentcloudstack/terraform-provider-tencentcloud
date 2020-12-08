---
subcategory: "Cloud Load Balancer(CLB)"
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
  health_check_port          = 200
  health_check_type          = "HTTP"
  health_check_http_code     = 2
  health_check_http_version  = "HTTP/1.0"
  health_check_http_method   = "GET"
}
```

TCP/UDP Listener with tcp health check

```hcl
resource "tencentcloud_clb_listener" "listener_tcp" {
  clb_id                     = tencentcloud_clb_instance.clb_basic.id
  listener_name              = "listener_tcp"
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  session_expire_time        = 30
  scheduler                  = "WRR"
  health_check_type          = "TCP"
  health_check_port          = 200
}
```

TCP/UDP Listener with http health check

```hcl
resource "tencentcloud_clb_listener" "listener_tcp" {
  clb_id                     = tencentcloud_clb_instance.clb_basic.id
  listener_name              = "listener_tcp"
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  session_expire_time        = 30
  scheduler                  = "WRR"
  health_check_type          = "HTTP"
  health_check_http_domain   = "www.tencent.com"
  health_check_http_code     = 16
  health_check_http_version  = "HTTP/1.1"
  health_check_http_method   = "HEAD"
  health_check_http_path     = "/"
}
```

TCP/UDP Listener with customer health check

```hcl
resource "tencentcloud_clb_listener" "listener_tcp" {
  clb_id                     = tencentcloud_clb_instance.clb_basic.id
  listener_name              = "listener_tcp"
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  session_expire_time        = 30
  scheduler                  = "WRR"
  health_check_type          = "CUSTOM"
  health_check_context_type  = "HEX"
  health_check_send_context  = "0123456789ABCDEF"
  health_check_recv_context  = "ABCD"
  target_type                = "TARGETGROUP"
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
  certificate_id       = "VjANRdz8"
  certificate_ca_id    = "VfqO4zkB"
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
  certificate_id             = "VjANRdz8"
  certificate_ca_id          = "VfqO4zkB"
  health_check_switch        = true
  health_check_time_out      = 2
  health_check_interval_time = 5
  health_check_health_num    = 3
  health_check_unhealth_num  = 3
  scheduler                  = "WRR"
  target_type                = "TARGETGROUP"
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required, ForceNew) Id of the CLB.
* `listener_name` - (Required) Name of the CLB listener, and available values can only be Chinese characters, English letters, numbers, underscore and hyphen '-'.
* `protocol` - (Required, ForceNew) Type of protocol within the listener. Valid values: 'TCP', 'UDP', 'HTTP', 'HTTPS' and 'TCP_SSL'.
* `certificate_ca_id` - (Optional) Id of the client certificate. NOTES: Only supports listeners of 'HTTPS' and 'TCP_SSL' protocol and must be set when the ssl mode is 'MUTUAL'.
* `certificate_id` - (Optional) Id of the server certificate. NOTES: Only supports listeners of 'HTTPS' and 'TCP_SSL' protocol and must be set when it is available.
* `certificate_ssl_mode` - (Optional) Type of certificate. Valid values: 'UNIDIRECTIONAL', 'MUTUAL'. NOTES: Only supports listeners of 'HTTPS' and 'TCP_SSL' protocol and must be set when it is available.
* `health_check_context_type` - (Optional) Health check protocol. When the value of `health_check_type` of the health check protocol is `CUSTOM`, this field is required, which represents the input format of the health check. Valid values: `HEX`, `TEXT`.
* `health_check_health_num` - (Optional) Health threshold of health check, and the default is 3. If a success result is returned for the health check for 3 consecutive times, the backend CVM is identified as healthy. The value range is 2-10. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.
* `health_check_http_code` - (Optional) HTTP health check code of TCP listener. When the value of `health_check_type` of the health check protocol is `HTTP`, this field is required. Valid values: 1, 2, 4, 8, 16. 1 means http_1xx, 2 means http_2xx, 4 means http_3xx, 8 means http_4xx, 16 means http_5xx.
* `health_check_http_domain` - (Optional) HTTP health check domain of TCP listener.
* `health_check_http_method` - (Optional) HTTP health check method of TCP listener. Valid values: `HEAD`, `GET`.
* `health_check_http_path` - (Optional) HTTP health check path of TCP listener.
* `health_check_http_version` - (Optional) The HTTP version of the backend service. When the value of `health_check_type` of the health check protocol is `HTTP`, this field is required. Valid values: `HTTP/1.0`, `HTTP/1.1`.
* `health_check_interval_time` - (Optional) Interval time of health check. Valid value ranges: (5~300) sec. and the default is 5 sec. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.
* `health_check_port` - (Optional) The health check port is the port of the backend service by default. Unless you want to specify a specific port, it is recommended to leave it blank. Only applicable to TCP/UDP listener.
* `health_check_recv_context` - (Optional) It represents the result returned by the health check. When the value of `health_check_type` of the health check protocol is `CUSTOM`, this field is required. Only ASCII visible characters are allowed and the maximum length is 500. When `health_check_context_type` value is `HEX`, the characters of SendContext and RecvContext can only be selected in `0123456789ABCDEF` and the length must be even digits.
* `health_check_send_context` - (Optional) It represents the content of the request sent by the health check. When the value of `health_check_type` of the health check protocol is `CUSTOM`, this field is required. Only visible ASCII characters are allowed and the maximum length is 500. When `health_check_context_type` value is `HEX`, the characters of SendContext and RecvContext can only be selected in `0123456789ABCDEF` and the length must be even digits.
* `health_check_switch` - (Optional) Indicates whether health check is enabled.
* `health_check_time_out` - (Optional) Response timeout of health check. Valid value ranges: (2~60) sec. Default is 2 sec. Response timeout needs to be less than check interval. NOTES: Only supports listeners of 'TCP','UDP','TCP_SSL' protocol.
* `health_check_type` - (Optional) Protocol used for health check. Valid values: `CUSTOM`, `TCP`, `HTTP`.
* `health_check_unhealth_num` - (Optional) Unhealthy threshold of health check, and the default is 3. If a success result is returned for the health check 3 consecutive times, the CVM is identified as unhealthy. The value range is 2-10. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.
* `port` - (Optional, ForceNew) Port of the CLB listener.
* `scheduler` - (Optional) Scheduling method of the CLB listener, and available values are 'WRR' and 'LEAST_CONN'. The default is 'WRR'. NOTES: The listener of HTTP and 'HTTPS' protocol additionally supports the 'IP Hash' method. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.
* `session_expire_time` - (Optional) Time of session persistence within the CLB listener. NOTES: Available when scheduler is specified as 'WRR', and not available when listener protocol is 'TCP_SSL'. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.
* `sni_switch` - (Optional, ForceNew) Indicates whether SNI is enabled, and only supported with protocol 'HTTPS'. If enabled, you can set a certificate for each rule in `tencentcloud_clb_listener_rule`, otherwise all rules have a certificate.
* `target_type` - (Optional, ForceNew) Backend target type. Valid values: `NODE`, `TARGETGROUP`. `NODE` means to bind ordinary nodes, `TARGETGROUP` means to bind target group. NOTES: TCP/UDP/TCP_SSL listener must configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `listener_id` - Id of this CLB listener.


## Import

CLB listener can be imported using the id (version >= 1.47.0), e.g.

```
$ terraform import tencentcloud_clb_listener.foo lb-7a0t6zqb#lbl-hh141sn9
```

