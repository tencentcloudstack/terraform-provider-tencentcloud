---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_layer4_listeners"
sidebar_current: "docs-tencentcloud-datasource-gaap_layer4_listeners"
description: |-
  Use this data source to query gaap layer4 listeners.
---

# tencentcloud_gaap_layer4_listeners

Use this data source to query gaap layer4 listeners.

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
resource "tencentcloud_gaap_layer4_listener" "foo" {
  protocol        = "TCP"
  name            = "ci-test-gaap-4-listener"
  port            = 80
  realserver_type = "IP"
  proxy_id        = "${tencentcloud_gaap_proxy.foo.id}"
  health_check    = true
  interval        = 5
  connect_timeout = 2
  realserver_bind_set {
    id   = "${tencentcloud_gaap_realserver.foo.id}"
    ip   = "${tencentcloud_gaap_realserver.foo.ip}"
    port = 80
  }
}
data "tencentcloud_gaap_layer4_listeners" "foo" {
  protocol    = "TCP"
  proxy_id    = "${tencentcloud_gaap_proxy.foo.id}"
  listener_id = "${tencentcloud_gaap_layer4_listener.foo.id}"
}
```

## Argument Reference

The following arguments are supported:

* `protocol` - (Required) Protocol of the layer4 listener to be queried, and the available values include `TCP` and `UDP`.
* `proxy_id` - (Required) ID of the GAAP proxy to be queried.
* `listener_id` - (Optional) ID of the layer4 listener to be queried.
* `listener_name` - (Optional) Name of the layer4 listener to be queried.
* `port` - (Optional) Port of the layer4 listener to be queried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `listeners` - An information list of layer4 listeners. Each element contains the following attributes:
  * `connect_timeout` - Timeout of the health check response.
  * `create_time` - Creation time of the layer4 listener.
  * `health_check` - Indicates whether health check is enable.
  * `id` - ID of the layer4 listener.
  * `interval` - Interval of the health check
  * `name` - Name of the layer4 listener.
  * `port` - Port of the layer4 listener.
  * `protocol` - Protocol of the layer4 listener.
  * `realserver_type` - Type of the realserver.
  * `scheduler` - Scheduling policy of the layer4 listener.
  * `status` - Status of the layer4 listener.


