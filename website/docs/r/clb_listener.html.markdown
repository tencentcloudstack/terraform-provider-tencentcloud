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

### HTTP Listener

```hcl
resource "tencentcloud_clb_listener" "HTTP_listener" {
  clb_id        = "lb-0lh5au7v"
  listener_name = "test_listener"
  port          = 80
  protocol      = "HTTP"
}
```

### TCP/UDP Listener

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

### TCP/UDP Listener with tcp health check

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

### TCP/UDP Listener with http health check

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

### TCP/UDP Listener with customer health check

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

### HTTPS Listener with sigle certificate

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

### HTTPS Listener with multi certificates

```hcl
resource "tencentcloud_clb_listener" "HTTPS_listener" {
  clb_id        = "lb-l6cp6jt4"
  listener_name = "test_listener"
  port          = "80"
  protocol      = "HTTPS"
  sni_switch    = true

  multi_cert_info {
    ssl_mode = "UNIDIRECTIONAL"
    cert_id_list = [
      "LCYouprI",
      "JVO1alRN"
    ]
  }
}
```

### TCP SSL Listener

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

### Port Range Listener

```hcl
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-listener-test"
}

resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  port                = 1
  end_port            = 6
  protocol            = "TCP"
  listener_name       = "listener_basic"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "NODE"
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required, String, ForceNew) ID of the CLB.
* `listener_name` - (Required, String) Name of the CLB listener, and available values can only be Chinese characters, English letters, numbers, underscore and hyphen '-'.
* `protocol` - (Required, String, ForceNew) Type of protocol within the listener. Valid values: `TCP`, `UDP`, `HTTP`, `HTTPS`, `TCP_SSL` and `QUIC`.
* `certificate_ca_id` - (Optional, String) ID of the client certificate. NOTES: Only supports listeners of `HTTPS` and `TCP_SSL` protocol and must be set when the ssl mode is `MUTUAL`.
* `certificate_id` - (Optional, String) ID of the server certificate. NOTES: Only supports listeners of `HTTPS` and `TCP_SSL` protocol and must be set when it is available.
* `certificate_ssl_mode` - (Optional, String) Type of certificate. Valid values: `UNIDIRECTIONAL`, `MUTUAL`. NOTES: Only supports listeners of `HTTPS` and `TCP_SSL` protocol and must be set when it is available.
* `end_port` - (Optional, Int, ForceNew) This parameter is used to specify the end port and is required when creating a port range listener. Only one member can be passed in when inputting the `Ports` parameter, which is used to specify the start port. If you want to try the port range feature, please [submit a ticket](https://console.cloud.tencent.com/workorder/category).
* `h2c_switch` - (Optional, Bool, ForceNew) Enable H2C switch for intranet HTTP listener.
* `health_check_context_type` - (Optional, String) Health check protocol. When the value of `health_check_type` of the health check protocol is `CUSTOM`, this field is required, which represents the input format of the health check. Valid values: `HEX`, `TEXT`.
* `health_check_health_num` - (Optional, Int) Health threshold of health check, and the default is `3`. If a success result is returned for the health check for 3 consecutive times, the backend CVM is identified as healthy. The value range is 2-10. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.
* `health_check_http_code` - (Optional, Int) HTTP health check code of TCP listener, Valid value ranges: [1~31]. When the value of `health_check_type` of the health check protocol is `HTTP`, this field is required. Valid values: `1`, `2`, `4`, `8`, `16`. `1` means http_1xx, `2` means http_2xx, `4` means http_3xx, `8` means http_4xx, `16` means http_5xx.If you want multiple return codes to indicate health, need to add the corresponding values.
* `health_check_http_domain` - (Optional, String) HTTP health check domain of TCP listener.
* `health_check_http_method` - (Optional, String) HTTP health check method of TCP listener. Valid values: `HEAD`, `GET`.
* `health_check_http_path` - (Optional, String) HTTP health check path of TCP listener.
* `health_check_http_version` - (Optional, String) The HTTP version of the backend service. When the value of `health_check_type` of the health check protocol is `HTTP`, this field is required. Valid values: `HTTP/1.0`, `HTTP/1.1`.
* `health_check_interval_time` - (Optional, Int) Interval time of health check. Valid value ranges: [2~300] sec. and the default is 5 sec. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in `tencentcloud_clb_listener_rule`.
* `health_check_port` - (Optional, Int) The health check port is the port of the backend service by default. Unless you want to specify a specific port, it is recommended to leave it blank. Only applicable to TCP/UDP listener.
* `health_check_recv_context` - (Optional, String) It represents the result returned by the health check. When the value of `health_check_type` of the health check protocol is `CUSTOM`, this field is required. Only ASCII visible characters are allowed and the maximum length is 500. When `health_check_context_type` value is `HEX`, the characters of SendContext and RecvContext can only be selected in `0123456789ABCDEF` and the length must be even digits.
* `health_check_send_context` - (Optional, String) It represents the content of the request sent by the health check. When the value of `health_check_type` of the health check protocol is `CUSTOM`, this field is required. Only visible ASCII characters are allowed and the maximum length is 500. When `health_check_context_type` value is `HEX`, the characters of SendContext and RecvContext can only be selected in `0123456789ABCDEF` and the length must be even digits.
* `health_check_switch` - (Optional, Bool) Indicates whether health check is enabled.
* `health_check_time_out` - (Optional, Int) Response timeout of health check. Valid value ranges: [2~60] sec. Default is 2 sec. Response timeout needs to be less than check interval. NOTES: Only supports listeners of `TCP`,`UDP`,`TCP_SSL` protocol.
* `health_check_type` - (Optional, String) Protocol used for health check. Valid values: `CUSTOM`, `TCP`, `HTTP`,`HTTPS`, `PING`, `GRPC`.
* `health_check_unhealth_num` - (Optional, Int) Unhealthy threshold of health check, and the default is `3`. If a success result is returned for the health check 3 consecutive times, the CVM is identified as unhealthy. The value range is [2-10]. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in `tencentcloud_clb_listener_rule`.
* `health_source_ip_type` - (Optional, Int) Specifies the type of health check source IP. `0` (default): CLB VIP. `1`: 100.64 IP range.
* `keepalive_enable` - (Optional, Int) Whether to enable a persistent connection. This parameter is applicable only to HTTP and HTTPS listeners. Valid values: 0 (disable; default value) and 1 (enable).
* `multi_cert_info` - (Optional, List) Certificate information. You can specify multiple server-side certificates with different algorithm types. This parameter is only applicable to HTTPS listeners with the SNI feature not enabled. Certificate and MultiCertInfo cannot be specified at the same time.
* `port` - (Optional, Int, ForceNew) Port of the CLB listener.
* `scheduler` - (Optional, String) Scheduling method of the CLB listener, and available values are 'WRR' and 'LEAST_CONN'. The default is 'WRR'. NOTES: The listener of `HTTP` and `HTTPS` protocol additionally supports the `IP Hash` method. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in `tencentcloud_clb_listener_rule`.
* `session_expire_time` - (Optional, Int) Time of session persistence within the CLB listener. NOTES: Available when scheduler is specified as `WRR`, and not available when listener protocol is `TCP_SSL`. NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, HTTP/HTTPS listener needs to be configured in `tencentcloud_clb_listener_rule`.
* `session_type` - (Optional, String) Session persistence type. Valid values: `NORMAL`: the default session persistence type; `QUIC_CID`: session persistence by QUIC connection ID. The `QUIC_CID` value can only be configured in UDP listeners. If this field is not specified, the default session persistence type will be used.
* `snat_enable` - (Optional, Bool) Whether to enable SNAT.
* `sni_switch` - (Optional, Bool, ForceNew) Indicates whether SNI is enabled, and only supported with protocol `HTTPS`. If enabled, you can set a certificate for each rule in `tencentcloud_clb_listener_rule`, otherwise all rules have a certificate.
* `target_type` - (Optional, String) Backend target type. Valid values: `NODE`, `TARGETGROUP`. `NODE` means to bind ordinary nodes, `TARGETGROUP` means to bind target group. NOTES: TCP/UDP/TCP_SSL listener must configuration, HTTP/HTTPS listener needs to be configured in tencentcloud_clb_listener_rule.

The `multi_cert_info` object supports the following:

* `cert_id_list` - (Required, Set) List of server certificate ID.
* `ssl_mode` - (Required, String) Authentication type. Values: UNIDIRECTIONAL (one-way authentication), MUTUAL (two-way authentication).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `listener_id` - ID of this CLB listener.


## Import

CLB listener can be imported using the id (version >= 1.47.0), e.g.

```
$ terraform import tencentcloud_clb_listener.foo lb-7a0t6zqb#lbl-hh141sn9
```

