---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_listener_rule"
sidebar_current: "docs-tencentcloud-resource-clb_listener_rule"
description: |-
  Provide a resource to create a CLB listener rule.
---

# tencentcloud_clb_listener_rule

Provide a resource to create a CLB listener rule.

## Example Usage

```hcl
resource "tencentcloud_clb_listener_forward_rule" "rule" {
  listener_id                = "lbl-hh141sn9"
  clb_id                     = "lb-k2zjp9lv"
  domain                     = "foo.net"
  url                        = "/bar"
  health_check_switch        = "0"
  health_check_interval_time = "5"
  health_check_health_num    = "3"
  health_check_unhealth_num  = "3"
  health_check_http_code     = "http_1xx"
  health_check_http_path     = "Default Path"
  health_check_http_domain   = "Default Domain"
  health_check_http_method   = "GET"
  certificate_server_id      = "my server certificate ID "
  certificate_ca_id          = "my client certificate ID"
  session_expire_time        = "0"
  schedule                   = "WRR"
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required) ID of CLB instance. 
* `domain` - (Required, ForceNew) Domain name of the forwarding rules
* `listener_id` - (Required, ForceNew) ID of CLB listener.
* `url` - (Required, ForceNew) Url of the forwarding rules
* `certificate_ca_id` - (Optional, ForceNew) Id of the client certificate.If not set, the content, key, name of client certificate must be set when SSLMode is 'mutual', only supported by listeners of protocol 'HTTPS'. 
* `certificate_id` - (Optional, ForceNew) Id of the server certificate.If not set, the content, key, name of server certificate must be set, only supported by listeners of protocol 'HTTPS'. 
* `certificate_ssl_mode` - (Optional, ForceNew) Type of SSL Mode. Available values are 'UNIDRECTIONAL', 'MUTUAL' 
* `health_check_health_num` - (Optional) Health threshold of health check, and the default is 3. If a success result is returned for the health check for 3 consecutive times, the backend CVM is identified as healthy. The value range is 2-10
* `health_check_http_code` - (Optional) Path of health check (applicable only to HTTP/HTTPS check methods).
* `health_check_http_domain` - (Optional) Domain name of health check (applicable only to HTTP/HTTPS check methods)
* `health_check_http_method` - (Optional) Methods of health check (applicable only to HTTP/HTTPS check methods). Available values include 'HEAD' and 'GET'.
* `health_check_http_path` - (Optional) Path of health check (applicable only to HTTP/HTTPS check methods). 
* `health_check_interval_time` - (Optional) Interval time of health check. The value range is 5-300 sec, and the default is 5 sec.
* `health_check_switch` - (Optional) Indicates whether health check is enabled.
* `health_check_unhealth_num` - (Optional) Health threshold of health check, and the default is 3. If a success result is returned for the health check for 3 consecutive times, the backend CVM is identified as healthy. The value range is 2-10
* `scheduler` - (Optional) Scheduling method of the CLB listener, and available values include 'WRR' and 'LEAST_CONN'. The defaule is 'WRR'.
* `session_expire_time` - (Optional) Time of session persistence within the CLB listener.


