---
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

* `name` - (Required) Name of the layer7 listener, the maximum length is 30.
* `port` - (Required, ForceNew) Port of the layer7 listener.
* `protocol` - (Required, ForceNew) Protocol of the layer7 listener, the available values include `HTTP` and `HTTPS`.
* `proxy_id` - (Required, ForceNew) ID of the GAAP proxy.
* `auth_type` - (Optional, ForceNew) Authentication type of the layer7 listener. `0` is one-way authentication and `1` is mutual authentication. NOTES: Only supports listeners of `HTTPS` protocol.
* `certificate_id` - (Optional) Certificate ID of the layer7 listener. NOTES: Only supports listeners of `HTTPS` protocol.
* `client_certificate_id` - (Optional, **Deprecated**) It has been deprecated from version 1.26.0. Set `client_certificate_ids` instead. ID of the client certificate. Set only when `auth_type` is specified as mutual authentication. NOTES: Only supports listeners of `HTTPS` protocol.
* `client_certificate_ids` - (Optional) ID list of the client certificate. Set only when `auth_type` is specified as mutual authentication. NOTES: Only supports listeners of `HTTPS` protocol.
* `forward_protocol` - (Optional, ForceNew) Protocol type of the forwarding, the available values include `HTTP` and `HTTPS`. NOTES: Only supports listeners of `HTTPS` protocol.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Creation time of the layer7 listener.
* `status` - Status of the layer7 listener.


## Import

GAAP layer7 listener can be imported using the id, e.g.

```
  $ terraform import tencentcloud_gaap_layer7_listener.foo listener-11112222
```

