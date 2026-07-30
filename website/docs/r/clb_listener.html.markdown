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
resource "tencentcloud_clb_listener" "example" {
  clb_id        = "lb-qck8thny"
  listener_name = "tf-example"
  port          = 80
  protocol      = "HTTP"
}
```

### TCP/UDP Listener

```hcl
resource "tencentcloud_clb_listener" "example" {
  clb_id                     = "lb-qck8thny"
  listener_name              = "tf-example"
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
  health_check_http_path     = "/"
  health_check_http_code     = 2
  health_check_http_version  = "HTTP/1.0"
  health_check_http_method   = "GET"
  deregister_target_rst      = false
  idle_connect_timeout       = 900
}
```

### TCP/UDP Listener with tcp health check

```hcl
resource "tencentcloud_clb_listener" "example" {
  clb_id                     = "lb-qck8thny"
  listener_name              = "tf-example"
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
  deregister_target_rst      = false
  idle_connect_timeout       = 900
}
```

### TCP/UDP Listener with http health check

```hcl
resource "tencentcloud_clb_listener" "example" {
  clb_id                     = "lb-qck8thny"
  listener_name              = "tf-example"
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
  deregister_target_rst      = false
  idle_connect_timeout       = 900
}
```

### TCP/UDP Listener with customer health check

```hcl
resource "tencentcloud_clb_listener" "example" {
  clb_id                     = "lb-qck8thny"
  listener_name              = "tf-example"
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
resource "tencentcloud_clb_listener" "example" {
  clb_id               = "lb-0lh5au7v"
  listener_name        = "tf-example"
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
resource "tencentcloud_clb_listener" "example" {
  clb_id        = "lb-l6cp6jt4"
  listener_name = "tf-example"
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
resource "tencentcloud_clb_listener" "example" {
  clb_id                     = "lb-0lh5au7v"
  listener_name              = "tf-example"
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

### TCP_SSL Listener with MaxConn, MaxCps, ProxyProtocol and DataCompressMode

```hcl
resource "tencentcloud_clb_listener" "example" {
  clb_id               = "lb-0lh5au7v"
  listener_name        = "tf-example"
  port                 = "443"
  protocol             = "TCP_SSL"
  certificate_ssl_mode = "UNIDIRECTIONAL"
  certificate_id       = "VjANRdz8"
  scheduler            = "WRR"
  max_conn             = 1000
  max_cps              = 100
  proxy_protocol       = true
  data_compress_mode   = "transparent"
}
```

### Port Range Listener

```hcl
resource "tencentcloud_clb_instance" "example" {
  clb_name     = "tf-listener-test"
  network_type = "OPEN"
}

resource "tencentcloud_clb_listener" "example" {
  clb_id              = tencentcloud_clb_instance.example.id
  listener_name       = "tf-example"
  port                = 1
  end_port            = 6
  protocol            = "TCP"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "NODE"
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required, String, ForceNew) ID of the CLB instance.
* `listener_name` - (Required, String) Name of the CLB listener, 1-80 characters. Supports letters, Chinese and other common international language characters, digits, hyphen '-' and underscore '_' (Unicode supplementary characters such as emoji are not allowed).
* `protocol` - (Required, String, ForceNew) Type of protocol within the listener. Valid values: `TCP`, `UDP`, `HTTP`, `HTTPS`, `TCP_SSL` and `QUIC`.
* `certificate_ca_id` - (Optional, String) ID of the client certificate. NOTES: Only supports listeners of `HTTPS` and `TCP_SSL` protocol and must be set when the ssl mode is `MUTUAL`.
* `certificate_id` - (Optional, String) ID of the server certificate. NOTES: Only supports listeners of `HTTPS` and `TCP_SSL` protocol and must be set when it is available.
* `certificate_ssl_mode` - (Optional, String) Type of certificate. Valid values: `UNIDIRECTIONAL`, `MUTUAL`. NOTES: Only supports listeners of `HTTPS` and `TCP_SSL` protocol and must be set when it is available.
* `data_compress_mode` - (Optional, String) Data compression mode. Valid values: `transparent`, `compatibility`.
* `deregister_target_rst` - (Optional, Bool) Reschedule function: the switch for unbinding backend services. When enabled, rescheduling is triggered when a backend service is unbound. Only supported by `TCP`/`UDP` listeners.
* `end_port` - (Optional, Int, ForceNew) This parameter is used to specify the end port and is required when creating a port range listener. Only one member can be passed in when inputting the `Ports` parameter, which is used to specify the start port. If you want to try the port range feature, please [submit a ticket](https://console.cloud.tencent.com/workorder/category).
* `h2c_switch` - (Optional, Bool, ForceNew) Whether to enable H2C for intranet `HTTP` listeners. `true`: enable, `false`: disable (default). When enabled, the listener only supports creating L7 rules with backend forwarding type `GRPC` or `GRPCS`; `GRPC` or `GRPCS` must be explicitly specified in the forwarding type when creating rules.
* `health_check_context_type` - (Optional, String) Custom probe parameter. Required when `health_check_type` is `CUSTOM`, representing the input format of the health check. Valid values: `HEX`, `TEXT`. When `HEX`, the characters of `send_context`/`recv_context` can only be selected from `0123456789ABCDEF` and the length must be even. Only applicable to `TCP`/`UDP` listeners.
* `health_check_health_num` - (Optional, Int) Health threshold. Default 3, meaning the backend is considered healthy after 3 consecutive successful probes. Value range: 2-10.
* `health_check_http_code` - (Optional, Int) Health check status code (only applicable to HTTP/HTTPS forwarding rules and the HTTP health check method of TCP listeners). Value range: 1-31, default 31. `1`=1xx healthy, `2`=2xx, `4`=3xx, `8`=4xx, `16`=5xx. To treat multiple return codes as healthy, add the corresponding values together.
* `health_check_http_domain` - (Optional, String) Health check domain, carried in the HTTP Host header (only applicable to HTTP/HTTPS listeners and the HTTP health check method of TCP listeners; for TCP listeners using HTTP health check, this field is required).
* `health_check_http_method` - (Optional, String) Health check method (only applicable to HTTP/HTTPS forwarding rules and the HTTP health check method of TCP listeners). Default `HEAD`. Valid values: `HEAD`, `GET`.
* `health_check_http_path` - (Optional, String) Health check path (only applicable to HTTP/HTTPS forwarding rules and the HTTP health check method of TCP listeners).
* `health_check_http_version` - (Optional, String) HTTP version of the backend service. Required when `health_check_type` is `HTTP`. Valid values: `HTTP/1.0`, `HTTP/1.1`. Only applicable to `TCP` listeners.
* `health_check_interval_time` - (Optional, Int) Health check probe interval in seconds. Default 5. Value range: 2-300 for IPv4 CLB instances and 5-300 for IPv6 CLB instances. Note: some older IPv4 CLB instances have a range of 5-300.
* `health_check_port` - (Optional, Int) Health check port. Defaults to the backend service port; leave blank unless a specific port is required. Pass `-1` to restore the default. Only applicable to `TCP`/`UDP` listeners.
* `health_check_recv_context` - (Optional, String) Custom probe parameter. Required when `health_check_type` is `CUSTOM`, representing the result returned by the health check. Only ASCII visible characters are allowed, max length 500. Only applicable to `TCP`/`UDP` listeners.
* `health_check_send_context` - (Optional, String) Custom probe parameter. Required when `health_check_type` is `CUSTOM`, representing the request content sent by the health check. Only ASCII visible characters are allowed, max length 500. Only applicable to `TCP`/`UDP` listeners.
* `health_check_switch` - (Optional, Bool) Indicates whether health check is enabled.
* `health_check_time_out` - (Optional, Int) Response timeout of health check in seconds. Value range: 2-60, default 2. The response timeout must be less than the check interval.
* `health_check_type` - (Optional, String) Health check protocol. Valid values: `TCP`, `HTTP`, `HTTPS`, `GRPC`, `PING`, `CUSTOM`. UDP listeners support `PING`/`CUSTOM`; TCP listeners support `TCP`/`HTTP`/`CUSTOM`; TCP_SSL/QUIC listeners support `TCP`/`HTTP`; HTTP rules support `HTTP`/`GRPC`; HTTPS rules support `HTTP`/`HTTPS`/`GRPC`. Defaults: `HTTP` for HTTP listeners, `TCP` for TCP/TCP_SSL/QUIC listeners, `PING` for UDP listeners; for HTTPS listeners the default matches the backend forwarding protocol.
* `health_check_unhealth_num` - (Optional, Int) Unhealthy threshold. Default 3, meaning the backend is considered unhealthy after 3 consecutive failed probes. Value range: 2-10.
* `health_source_ip_type` - (Optional, Int) Health check source IP type. `0`: use the CLB VIP as the source IP, `1`: use a 100.64 IP range as the source IP.
* `idle_connect_timeout` - (Optional, Int) Idle connection timeout. This parameter is only available for TCP/UDP listeners, in seconds. Default: 900s for TCP listeners, 300s for UDP listeners. Value range: 10-900 for shared and dedicated instances; 10-1980 for LCU-supported CLB instances. To set a value beyond the range, please submit a ticket for application.
* `keepalive_enable` - (Optional, Int) Whether to enable persistent connection (long connection). Only applicable to `HTTP`/`HTTPS` listeners. Valid values: `0` (disable, default), `1` (enable). This feature is currently in beta.
* `max_conn` - (Optional, Int) Listener-level maximum concurrent connections. Currently only supported for performance capacity-type CLB instances with TCP/UDP/TCP_SSL/QUIC listeners. Pass -1 to indicate no limit at the listener level. Basic network instances do not support this parameter.
* `max_cps` - (Optional, Int) Listener-level maximum new connections per second. Currently only supported for performance capacity-type CLB instances with TCP/UDP/TCP_SSL/QUIC listeners. Pass -1 to indicate no limit at the listener level. Basic network instances do not support this parameter.
* `multi_cert_info` - (Optional, List) Certificate information, supporting multiple server certificates with different algorithm types at the same time. Only applicable to `TCP_SSL` listeners and `HTTPS` listeners with SNI disabled. When creating a `TCP_SSL` listener or an `HTTPS` listener with SNI disabled, at least one of `certificate`/`multi_cert_info` must be specified, but they cannot be specified at the same time.
* `port` - (Optional, Int, ForceNew) Port of the CLB listener. Port range: [1 - 65535].
* `proxy_protocol` - (Optional, Bool) Enable proxy protocol for TCP_SSL and QUIC listeners. Note: this field is not returned by the DescribeListeners API, so it will not be refreshed in state after creation.
* `reschedule_expand_target` - (Optional, Bool) The rescheduling function, a switch for scaling backend services, triggers rescheduling when backend servers are added or removed. Only supported by TCP/UDP listeners.
* `reschedule_interval` - (Optional, Int) Rescheduled trigger duration, ranging from 0 to 3600 seconds. Supported only by TCP/UDP listeners.
* `reschedule_start_time` - (Optional, Int) Reschedule the trigger start time, with a value ranging from 0 to 3600 seconds. Only supported by TCP/UDP listeners.
* `reschedule_target_zero_weight` - (Optional, Bool) The rescheduling function, with a weight of 0 as a switch, triggers rescheduling when the weight of the backend server is set to 0. Only supported by TCP/UDP listeners.
* `reschedule_unhealthy` - (Optional, Bool) Rescheduling function, health check exception switch. Enabling this switch triggers rescheduling when a backend server fails a health check. Supported only by TCP/UDP listeners.
* `scheduler` - (Optional, String) Scheduling method. Valid values: `WRR` (weighted round-robin), `LEAST_CONN` (least connections). Default is `WRR`. Only applicable to `TCP`/`UDP`/`TCP_SSL`/`QUIC` listeners.
* `session_expire_time` - (Optional, Int) Session persistence time in seconds. Value range: 30-3600, default 0 (disabled). Only applicable to `TCP`/`UDP` listeners.
* `session_type` - (Optional, String) Session persistence type. `NORMAL` (default): default session persistence type; `QUIC_CID`: session persistence by QUIC connection ID. Only applicable to `TCP`/`UDP` listeners; L7 listeners should be configured in the forwarding rule. If `QUIC_CID` is selected, `protocol` must be `UDP`, `scheduler` must be `WRR`, and only IPv4 is supported.
* `snat_enable` - (Optional, Bool) Whether to enable SNAT (source IP replacement). `true`: enable, `false`: disable (default). Note: when SNAT is enabled, the client source IP is replaced and the pass-through client source IP option is disabled, and vice versa.
* `sni_switch` - (Optional, Bool, ForceNew) Indicates whether SNI is enabled. Only applicable to `HTTPS` listeners. `0`: disabled, `1`: enabled.
* `target_type` - (Optional, String) Backend target type. Valid values: `NODE`, `TARGETGROUP`, `TARGETGROUP-V2`. `NODE` means binding ordinary nodes, `TARGETGROUP` means binding a target group. Only applicable to `TCP`/`UDP` listeners; L7 listeners should be configured in the forwarding rule.

The `multi_cert_info` object supports the following:

* `cert_id_list` - (Required, Set) List of server certificate ID.
* `ssl_mode` - (Required, String) Authentication type. Values: UNIDIRECTIONAL (one-way authentication), MUTUAL (two-way authentication).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `listener_id` - ID of this CLB listener.


## Import

CLB listener can be imported using the clbId#listenerId (version >= 1.47.0), e.g.

```
$ terraform import tencentcloud_clb_listener.example lb-7a0t6zqb#lbl-hh141sn9
```

