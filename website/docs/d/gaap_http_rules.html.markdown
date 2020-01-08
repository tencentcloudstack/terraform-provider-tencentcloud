---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_http_rules"
sidebar_current: "docs-tencentcloud-datasource-gaap_http_rules"
description: |-
  Use this data source to query forward rule of layer7 listeners.
---

# tencentcloud_gaap_http_rules

Use this data source to query forward rule of layer7 listeners.

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

resource "tencentcloud_gaap_realserver" "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource "tencentcloud_gaap_http_rule" "foo" {
  listener_id     = tencentcloud_gaap_layer7_listener.foo.id
  domain          = "www.qq.com"
  path            = "/"
  realserver_type = "IP"
  health_check    = true

  realservers {
    id   = tencentcloud_gaap_realserver.foo.id
    ip   = tencentcloud_gaap_realserver.foo.ip
    port = 80
  }
}

data "tencentcloud_gaap_http_rules" "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = tencentcloud_gaap_http_rule.foo.domain
}
```

## Argument Reference

The following arguments are supported:

* `listener_id` - (Required) ID of the layer7 listener to be queried.
* `domain` - (Optional) Forward domain of the layer7 listener to be queried.
* `forward_host` - (Optional) Requested host which is forwarded to the realserver by the listener to be queried.
* `path` - (Optional) Path of the forward rule to be queried.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `rules` - An information list of forward rule of the layer7 listeners. Each element contains the following attributes:
  * `connect_timeout` - Timeout of the health check response.
  * `domain` - Forward domain of the forward rule.
  * `forward_host` - Requested host which is forwarded to the realserver by the listener.
  * `health_check_method` - Method of the health check.
  * `health_check_path` - Path of health check.
  * `health_check_status_codes` - Return code of confirmed normal.
  * `health_check` - Indicates whether health check is enable.
  * `id` - ID of the forward rule.
  * `interval` - Interval of the health check.
  * `listener_id` - ID of the layer7 listener.
  * `path` - Path of the forward rule.
  * `realserver_type` - Type of the realserver.
  * `realservers` - An information list of GAAP realserver. Each element contains the following attributes:
    * `domain` - Domain of the GAAP realserver.
    * `id` - ID of the GAAP realserver.
    * `ip` - IP of the GAAP realserver.
    * `port` - Port of the GAAP realserver.
    * `status` - Status of the GAAP realserver.
    * `weight` - Scheduling weight.
  * `scheduler` - Scheduling policy of the forward rule.


