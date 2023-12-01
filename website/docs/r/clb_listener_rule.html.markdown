---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_listener_rule"
sidebar_current: "docs-tencentcloud-resource-clb_listener_rule"
description: |-
  Provides a resource to create a CLB listener rule.
---

# tencentcloud_clb_listener_rule

Provides a resource to create a CLB listener rule.

-> **NOTE:** This resource only be applied to the HTTP or HTTPS listeners.

## Example Usage

```hcl
resource "tencentcloud_clb_listener_rule" "foo" {
  listener_id                = "lbl-hh141sn9"
  clb_id                     = "lb-k2zjp9lv"
  domain                     = "foo.net"
  url                        = "/bar"
  health_check_switch        = true
  health_check_interval_time = 5
  health_check_health_num    = 3
  health_check_unhealth_num  = 3
  health_check_http_code     = 2
  health_check_http_path     = "Default Path"
  health_check_http_domain   = "Default Domain"
  health_check_http_method   = "GET"
  certificate_ssl_mode       = "MUTUAL"
  certificate_id             = "VjANRdz8"
  certificate_ca_id          = "VfqO4zkB"
  session_expire_time        = 30
  scheduler                  = "WRR"
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required, String) ID of CLB instance.
* `domain` - (Required, String) Domain name of the listener rule.
* `listener_id` - (Required, String, ForceNew) ID of CLB listener.
* `url` - (Required, String) Url of the listener rule.
* `certificate_ca_id` - (Optional, String) ID of the client certificate. NOTES: Only supports listeners of HTTPS protocol.
* `certificate_id` - (Optional, String) ID of the server certificate. NOTES: Only supports listeners of HTTPS protocol.
* `certificate_ssl_mode` - (Optional, String, ForceNew) Type of certificate. Valid values: `UNIDIRECTIONAL`, `MUTUAL`. NOTES: Only supports listeners of HTTPS protocol.
* `forward_type` - (Optional, String) Forwarding protocol between the CLB instance and real server. Valid values: `HTTP`, `HTTPS`, `TRPC`. The default is `HTTP`.
* `health_check_health_num` - (Optional, Int) Health threshold of health check, and the default is `3`. If a success result is returned for the health check 3 consecutive times, indicates that the forwarding is normal. The value range is [2-10]. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in `tencentcloud_clb_listener_rule`.
* `health_check_http_code` - (Optional, Int) HTTP Status Code. The default is 31. Valid value ranges: [1~31]. `1 means the return value '1xx' is health. `2` means the return value '2xx' is health. `4` means the return value '3xx' is health. `8` means the return value '4xx' is health. 16 means the return value '5xx' is health. If you want multiple return codes to indicate health, need to add the corresponding values. NOTES: The 'HTTP' health check of the 'TCP' listener only supports specifying one health check status code. NOTES: Only supports listeners of 'HTTP' and 'HTTPS' protocol.
* `health_check_http_domain` - (Optional, String) Domain name of health check. NOTES: Only supports listeners of `HTTP` and `HTTPS` protocol.
* `health_check_http_method` - (Optional, String) Methods of health check. NOTES: Only supports listeners of `HTTP` and `HTTPS` protocol. The default is `HEAD`, the available value are `HEAD` and `GET`.
* `health_check_http_path` - (Optional, String) Path of health check. NOTES: Only supports listeners of `HTTP` and `HTTPS` protocol.
* `health_check_interval_time` - (Optional, Int) Interval time of health check. Valid value ranges: (2~300) sec. and the default is `5` sec. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in `tencentcloud_clb_listener_rule`.
* `health_check_switch` - (Optional, Bool) Indicates whether health check is enabled.
* `health_check_time_out` - (Optional, Int) Time out of health check. The value range is [2-60](SEC).
* `health_check_type` - (Optional, String) Type of health check. Valid value is `CUSTOM`, `TCP`, `HTTP`.
* `health_check_unhealth_num` - (Optional, Int) Unhealthy threshold of health check, and the default is `3`. If the unhealthy result is returned 3 consecutive times, indicates that the forwarding is abnormal. The value range is [2-10].  NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in `tencentcloud_clb_listener_rule`.
* `http2_switch` - (Optional, Bool) Indicate to apply HTTP2.0 protocol or not.
* `scheduler` - (Optional, String) Scheduling method of the CLB listener rules. Valid values: `WRR`, `IP HASH`, `LEAST_CONN`. The default is `WRR`.  NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in `tencentcloud_clb_listener_rule`.
* `session_expire_time` - (Optional, Int) Time of session persistence within the CLB listener. NOTES: Available when scheduler is specified as `WRR`, and not available when listener protocol is `TCP_SSL`.  NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in `tencentcloud_clb_listener_rule`.
* `target_type` - (Optional, String, ForceNew) Backend target type. Valid values: `NODE`, `TARGETGROUP`. `NODE` means to bind ordinary nodes, `TARGETGROUP` means to bind target group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - ID of this CLB listener rule.


## Import

CLB listener rule can be imported using the id (version >= 1.47.0), e.g.

```
$ terraform import tencentcloud_clb_listener_rule.foo lb-7a0t6zqb#lbl-hh141sn9#loc-agg236ys
```

