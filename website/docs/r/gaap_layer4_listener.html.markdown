---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_layer4_listener"
sidebar_current: "docs-tencentcloud-resource-gaap_layer4_listener"
description: |-
  Provides a resource to create a layer4 listener of GAAP.
---

# tencentcloud_gaap_layer4_listener

Provides a resource to create a layer4 listener of GAAP.

## Example Usage

```hcl
resource "tencentcloud_gaap_proxy" "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource "tencentcloud_gaap_realserver" "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource "tencentcloud_gaap_realserver" "bar" {
  ip   = "119.29.29.29"
  name = "ci-test-gaap-realserver2"
}

resource "tencentcloud_gaap_layer4_listener" "foo" {
  protocol        = "TCP"
  name            = "ci-test-gaap-4-listener"
  port            = 80
  realserver_type = "IP"
  proxy_id        = "${tencentcloud_gaap_proxy.foo.id}"
  health_check    = true

  realserver_bind_set {
    id   = "${tencentcloud_gaap_realserver.foo.id}"
    ip   = "${tencentcloud_gaap_realserver.foo.ip}"
    port = 80
  }

  realserver_bind_set {
    id   = "${tencentcloud_gaap_realserver.bar.id}"
    ip   = "${tencentcloud_gaap_realserver.bar.ip}"
    port = 80
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the layer4 listener, the maximum length is 30.
* `port` - (Required, ForceNew) Port of the layer4 listener.
* `protocol` - (Required, ForceNew) Protocol of the layer4 listener, and the available values include `TCP` and `UDP`.
* `proxy_id` - (Required, ForceNew) ID of the GAAP proxy.
* `realserver_type` - (Required, ForceNew) Type of the realserver, and the available values include `IP`,`DOMAIN`. NOTES: when the `protocol` is specified as `TCP` and the `scheduler` is specified as `wrr`, the item can only be set to `IP`.
* `connect_timeout` - (Optional) Timeout of the health check response, default is 2s. NOTES: Only supports listeners of `TCP` protocol and require less than `interval`.
* `health_check` - (Optional) Indicates whether health check is enable, default is false. NOTES: Only supports listeners of `TCP` protocol.
* `interval` - (Optional) Interval of the health check, default is 5s. NOTES: Only supports listeners of `TCP` protocol.
* `realserver_bind_set` - (Optional) An information list of GAAP realserver. Each element contains the following attributes:
* `scheduler` - (Optional) Scheduling policy of the layer4 listener, default is `rr`. Available values include `rr`,`wrr` and `lc`.

The `realserver_bind_set` object supports the following:

* `id` - (Required) ID of the GAAP realserver.
* `ip` - (Required) IP of the GAAP realserver.
* `port` - (Required) Port of the GAAP realserver.
* `weight` - (Optional) Scheduling weight, default is 1. The range of values is [1,100].

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Creation time of the layer4 listener.
* `status` - Status of the layer4 listener.


