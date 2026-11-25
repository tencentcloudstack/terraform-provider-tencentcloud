---
subcategory: "Application Load Balancer(alb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_alb_listener"
sidebar_current: "docs-tencentcloud-resource-alb_listener"
description: |-
  Provides a resource to manage an ALB listener.
---

# tencentcloud_alb_listener

Provides a resource to manage an ALB listener.

## Example Usage

```hcl
resource "tencentcloud_alb_listener" "example" {
  load_balancer_id  = "alb-12345678"
  listener_protocol = "HTTPS"
  listener_port     = 443
  listener_name     = "tf-example-listener"

  default_actions {
    type = "TargetGroup"
    target_group_config {
      target_groups {
        target_group_id = "lbtg-12345678"
        weight          = 10
      }
    }
  }

  tags {
    tag_key   = "env"
    tag_value = "tf-example"
  }
}
```

## Argument Reference

The following arguments are supported:

* `default_actions` - (Required, List) List of default forwarding rule actions. Currently, only 1 default forwarding rule action can be added to a listener.
* `listener_port` - (Required, Int) Port used by the load balancer instance frontend. Value range: 1~65535.
* `listener_protocol` - (Required, String) Listener protocol. Valid values: `HTTP`, `HTTPS`.
* `load_balancer_id` - (Required, String) Load balancer instance ID, in the format of `alb-` followed by 8 alphanumeric characters.
* `ca_certificate_ids` - (Optional, List) List of CA certificate IDs configured for the listener. Currently, only 1 CA certificate can be added to a listener. Required when `ca_enabled` is set to `true`.
* `ca_enabled` - (Optional, Bool) Whether to enable mutual authentication. Valid values: `true`: enabled; `false` (default): disabled. Read back from the cloud as `CaEnable`.
* `certificate_ids` - (Optional, List) List of server certificate IDs.
* `client_token` - (Optional, String) Client token used to guarantee request idempotency. Generate a parameter value from your client ensuring the parameter value is unique across different requests. Only ASCII characters are supported.
* `gzip_enabled` - (Optional, Bool) Whether to enable Gzip compression. Valid values: `true` (default): yes; `false`: no.
* `http2_enabled` - (Optional, Bool) Whether to enable the HTTP/2 feature. `false` by default for the HTTP protocol and `true` by default for the HTTPS protocol. Only the HTTPS protocol supports this parameter. Read back from the cloud as `Http2Enable`.
* `idle_timeout` - (Optional, Int) Connection idle timeout in seconds. Value range: 1~600. Default value: 15. If no access request arrives within the timeout period, the load balancer disconnects the current connection and creates a new one when the next request arrives.
* `listener_name` - (Optional, String) Custom listener name. 1~255 characters long, must be harmless characters, and can contain Chinese characters, letters, digits, hyphen (-), forward slash (/), dot (.) and underscore (_).
* `request_timeout` - (Optional, Int) Connection request timeout in seconds. Value range: 1~600. Default value: 60. If the backend server does not return a response within the timeout period, the load balancer gives up waiting and returns an HTTP 504 error to the client.
* `security_policy_id` - (Optional, String) Security policy ID, in the format of `tls-` followed by 8 alphanumeric characters. Read back from the cloud as `TlsSecurityPolicyId`.
* `tags` - (Optional, List) List of tags. Up to 20 supported.
* `x_forwarded_for_config` - (Optional, Object) X-Forwarded-For configuration.

The `default_actions` object supports the following:

* `target_group_config` - (Required, Object) Target group configuration of the forwarding action. When creating a listener, the target group configuration in the forwarding action only supports a single target group.
* `type` - (Required, String) Forwarding action type. When creating a listener, the default forwarding action type only supports forwarding to a target group.

The `target_group_config` object of `default_actions` supports the following:

* `target_groups` - (Required, List) List of target groups.
* `target_group_sticky_session` - (Optional, Object) Sticky session between target groups.

The `target_group_sticky_session` object of `target_group_config` supports the following:

* `enabled` - (Required, Bool) Whether to enable sticky session. Disabled by default.
* `timeout` - (Optional, Int) Timeout in seconds. Value range: 1~86400. Default value: 1000.

The `target_groups` object of `target_group_config` supports the following:

* `target_group_id` - (Required, String) Target group ID, in the format of `lbtg-` followed by 8 alphanumeric characters.
* `weight` - (Optional, Int) Weight. Value range: [0, 100]. Default value: 10.

The `tags` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

The `x_forwarded_for_config` object supports the following:

* `x_forwarded_for_alb_id_enabled` - (Optional, Bool) Whether to obtain the load balancer instance ID through the ALB-ID header field. `true`: yes; `false`: no.
* `x_forwarded_for_client_src_port_enabled` - (Optional, Bool) Whether to obtain the port of the client accessing the load balancer instance through the X-Forwarded-Client-srcport header field. `true`: yes; `false`: no.
* `x_forwarded_for_host_enabled` - (Optional, Bool) Whether to obtain the domain name of the client accessing the load balancer instance through the X-Forwarded-Host header field. `true`: yes; `false`: no.
* `x_forwarded_for_mode` - (Optional, String) Specifies how to process the X-Forwarded-For (XFF) HTTP header field. Valid values: `append`: append mode (default), which appends the client's real IP address to the end of the X-Forwarded-For header, retaining the original XFF chain information; `remove`: remove mode, which removes the X-Forwarded-For header and does not pass it to the backend server; `passthrough`: pass-through mode, which keeps the X-Forwarded-For header unchanged and passes it to the backend server without any modification.
* `x_forwarded_for_port_enabled` - (Optional, Bool) Whether to obtain the listener port of the load balancer instance through the X-Forwarded-Port header field. `true`: yes; `false`: no.
* `x_forwarded_for_proto_enabled` - (Optional, Bool) Whether to obtain the listener protocol of the load balancer instance through the X-Forwarded-Proto header field. `true`: yes; `false`: no.
* `x_tencent_client_idn_enabled` - (Optional, Bool) Whether to access the issuer of the client certificate, $ssl_client_i_dn, through the X-Tencent-Client-IDN header. `true`: yes; `false`: no.
* `x_tencent_client_sdn_enabled` - (Optional, Bool) Whether to access the subject of the client certificate, $ssl_client_s_dn, through the X-Tencent-Client-SDN header. `true`: yes; `false`: no.
* `x_tencent_client_serial_enabled` - (Optional, Bool) Whether to access the serial number of the client certificate, $ssl_client_serial, through the X-Tencent-Client-Serial header. `true`: yes; `false`: no.
* `x_tencent_client_verify_enabled` - (Optional, Bool) Whether to access the verification result of the client certificate, $ssl_client_verify, through the X-Tencent-Client-Verify header. `true`: yes; `false`: no.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource, the listener ID, in the format of `lst-` followed by 8 alphanumeric characters.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateListener](https://www.tencentcloud.com/document/api/1311/84262)
  - [DescribeListeners](https://www.tencentcloud.com/document/api/1311/84257)
* Read: [DescribeListeners](https://www.tencentcloud.com/document/api/1311/84257)
* Update:
  - [DescribeListeners](https://www.tencentcloud.com/document/api/1311/84257)
  - [ModifyListenerAttributes](https://www.tencentcloud.com/document/api/1311/84255)
* Delete: [DeleteListener](https://www.tencentcloud.com/document/api/1311/84261)

## Import

ALB listener can be imported using the id, e.g.

```
terraform import tencentcloud_alb_listener.example lst-12345678
```

