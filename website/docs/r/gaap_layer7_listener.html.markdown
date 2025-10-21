---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_layer7_listener"
sidebar_current: "docs-tencentcloud-resource-gaap_layer7_listener"
description: |-
  Provides a resource to create a layer7 listener of GAAP.
---

# tencentcloud_gaap_layer7_listener

Provides a resource to create a layer7 listener of GAAP.

## Example Usage

```hcl
resource "tencentcloud_gaap_proxy" "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource "tencentcloud_gaap_layer7_listener" "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = tencentcloud_gaap_proxy.foo.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Name of the layer7 listener, the maximum length is 30.
* `port` - (Required, Int, ForceNew) Port of the layer7 listener.
* `protocol` - (Required, String, ForceNew) Protocol of the layer7 listener. Valid value: `HTTP` and `HTTPS`.
* `auth_type` - (Optional, Int, ForceNew) Authentication type of the layer7 listener. `0` is one-way authentication and `1` is mutual authentication. NOTES: Only supports listeners of `HTTPS` protocol.
* `certificate_id` - (Optional, String) Certificate ID of the layer7 listener. NOTES: Only supports listeners of `HTTPS` protocol.
* `client_certificate_id` - (Optional, String, **Deprecated**) It has been deprecated from version 1.26.0. Set `client_certificate_ids` instead. ID of the client certificate. Set only when `auth_type` is specified as mutual authentication. NOTES: Only supports listeners of `HTTPS` protocol.
* `client_certificate_ids` - (Optional, Set: [`String`]) ID list of the client certificate. Set only when `auth_type` is specified as mutual authentication. NOTES: Only supports listeners of `HTTPS` protocol.
* `forward_protocol` - (Optional, String, ForceNew) Protocol type of the forwarding. Valid value: `HTTP` and `HTTPS`. NOTES: Only supports listeners of `HTTPS` protocol.
* `group_id` - (Optional, String, ForceNew) Group ID.
* `proxy_id` - (Optional, String, ForceNew) ID of the GAAP proxy.
* `tls_ciphers` - (Optional, String) Password Suite, optional GAAP_TLS_CIPHERS_STRICT, GAAP_TLS_CIPHERS_GENERAL, GAAP_TLS_CIPHERS_WIDE(default).
* `tls_support_versions` - (Optional, Set: [`String`]) TLS version, optional TLSv1, TLSv1.1, TLSv1.2, TLSv1.3.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of the layer7 listener.
* `status` - Status of the layer7 listener.


## Import

GAAP layer7 listener can be imported using the id, e.g.

```
  $ terraform import tencentcloud_gaap_layer7_listener.foo listener-11112222
```

