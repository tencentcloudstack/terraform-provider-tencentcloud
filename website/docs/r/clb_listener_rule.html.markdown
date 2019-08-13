---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_listener_rule"
sidebar_current: "docs-tencentcloud-resource-clb_listener_rule"
description: |-
  Provides a resource to create a CLB listener rule.
---

# tencentcloud_clb_listener_rule

Provides a resource to create a CLB listener rule.

~> **NOTE:** This resource only be created in the HTTP or HTTPS listeners.

## Example Usage

```hcl
resource "tencentcloud_clb_listener_rule" "foo" {
  listener_id                = "lbl-hh141sn9#lb-k2zjp9lv"
  clb_id                     = "lb-k2zjp9lv"
  domain                     = "foo.net"
  url                        = "/bar"
  health_check_switch        = true
  health_check_interval_time = 5
  health_check_health_num    = 3
  health_check_unhealth_num  = 3
  health_check_http_code     = "http_1xx"
  health_check_http_path     = "Default Path"
  health_check_http_domain   = "Default Domain"
  health_check_http_method   = "GET"
  certificate_ssl_mode       = "MUTUAL"
  certificate_id             = "mycert server ID "
  certificate_ca_id          = "mycert ca ID"
  session_expire_time        = 30
  scheduler                  = "WRR"
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required) ID of CLB instance.
* `domain` - (Required, ForceNew) Domain name of the forwarding rule.
* `listener_id` - (Required, ForceNew) ID of CLB listener.
* `url` - (Required, ForceNew) Url of the forwarding rule.
* `certificate_ca_id` - (Optional, ForceNew) ID of the client certificate. NOTES: Only supports listeners of 'HTTPS' and 'TCP_SSL' protocol.
* `certificate_id` - (Optional, ForceNew) ID of the server certificate. NOTES: Only supports listeners of 'HTTPS' and 'TCP_SSL' protocol.
* `certificate_ssl_mode` - (Optional, ForceNew) Type of certificate, and available values inclue 'UNIDIRECTIONAL', 'MUTUAL'. NOTES: Only supports listeners of 'HTTPS' and 'TCP_SSL' protocol.
* `health_check_health_num` - (Optional) Health threshold of health check, and the default is 3. If a success result is returned for the health check 3 consecutive times, indicates that the forwarding is normal. The value range is 2-10.
* `health_check_http_code` - (Optional) HTTP Status Code. The default is 31 and value range is 1-31. '0b0001' means the return value '1xx' is health. '0b0010' means the return value '2xx' is health. '0b0100' means the return value '3xx' is health. '0b1000' means the return value '4xx' is health. 0b10000 means the return value '5xx' is health. If you want multiple return codes to indicate health, need to add the corresponding values. NOTES: The 'HTTP' health check of the 'TCP' listener only supports specifying one health check status code. NOTES: Only supports listeners of 'HTTP' and 'HTTPS' protocol.
* `health_check_http_domain` - (Optional) Domain name of health check. NOTES: Only supports listeners of 'HTTP' and 'HTTPS' protocol.
* `health_check_http_method` - (Optional) Methods of health check. NOTES: Only supports listeners of 'HTTP' and 'HTTPS' protocol. The default is 'HEAD', the available value include 'HEAD' and 'GET'.
* `health_check_http_path` - (Optional) Path of health check. NOTES: Only supports listeners of 'HTTP' and 'HTTPS' protocol.
* `health_check_interval_time` - (Optional) Interval time of health check. The value range is 5-300 sec, and the default is 5 sec.
* `health_check_switch` - (Optional) Indicates whether health check is enabled.
* `health_check_unhealth_num` - (Optional) Unhealth threshold of health check, and the default is 3. If the unhealth result is returned 3 consecutive times, indicates that the forwarding is abnormal. The value range is 2-10.
* `scheduler` - (Optional) Scheduling method of the CLB listener, and available values include 'WRR', 'IP HASH' and 'LEAST_CONN'. The default is 'WRR'.
* `session_expire_time` - (Optional) Time of session persistence within the CLB listener. NOTES: Available when scheduler is specified as 'WRR', and not available when listener protocol is TCP_SSL.


## Import

CLB instance rule can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_listener_rule.foo #loc-4xxr2cy7#lbl-hh141sn9#lb-k2zjp9lv
```

