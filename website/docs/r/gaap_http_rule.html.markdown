---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_http_rule"
sidebar_current: "docs-tencentcloud-resource-gaap_http_rule"
description: |-
  Provides a resource to create a forward rule of layer7 listener.
---

# tencentcloud_gaap_http_rule

Provides a resource to create a forward rule of layer7 listener.

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

resource "tencentcloud_gaap_realserver" "bar" {
  ip   = "8.8.8.8"
  name = "ci-test-gaap-realserver"
}

resource "tencentcloud_gaap_http_domain" "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"
}

resource "tencentcloud_gaap_http_rule" "foo" {
  listener_id               = tencentcloud_gaap_layer7_listener.foo.id
  domain                    = tencentcloud_gaap_http_domain.foo.domain
  path                      = "/"
  realserver_type           = "IP"
  health_check              = true
  health_check_path         = "/"
  health_check_method       = "GET"
  health_check_status_codes = [200]

  realservers {
    id   = tencentcloud_gaap_realserver.foo.id
    ip   = tencentcloud_gaap_realserver.foo.ip
    port = 80
  }

  realservers {
    id   = tencentcloud_gaap_realserver.bar.id
    ip   = tencentcloud_gaap_realserver.bar.ip
    port = 80
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, ForceNew) Forward domain of the forward rule.
* `health_check` - (Required) Indicates whether health check is enable.
* `listener_id` - (Required, ForceNew) ID of the layer7 listener.
* `path` - (Required) Path of the forward rule. Maximum length is 80.
* `realserver_type` - (Required, ForceNew) Type of the realserver. Valid value: `IP` and `DOMAIN`.
* `connect_timeout` - (Optional) Timeout of the health check response, default value is 2s.
* `forward_host` - (Optional) The default value of requested host which is forwarded to the realserver by the listener is `default`.
* `health_check_method` - (Optional) Method of the health check. Valid value: `GET` and `HEAD`.
* `health_check_path` - (Optional) Path of health check. Maximum length is 80.
* `health_check_status_codes` - (Optional) Return code of confirmed normal. Valid value: `100`, `200`, `300`, `400` and `500`.
* `interval` - (Optional) Interval of the health check, default value is 5s.
* `realservers` - (Optional) An information list of GAAP realserver.
* `scheduler` - (Optional) Scheduling policy of the forward rule, default value is `rr`. Valid value: `rr`, `wrr` and `lc`.
* `sni_switch` - (Optional) ServerNameIndication (SNI) switch. ON means on and OFF means off.
* `sni` - (Optional) ServerNameIndication (SNI) is required when the SNI switch is turned on.

The `realservers` object supports the following:

* `id` - (Required) ID of the GAAP realserver.
* `ip` - (Required) IP of the GAAP realserver.
* `port` - (Required) Port of the GAAP realserver.
* `weight` - (Optional) Scheduling weight, default value is `1`. Valid value ranges: (1~100).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

GAAP http rule can be imported using the id, e.g.

```
  $ terraform import tencentcloud_gaap_http_rule.foo rule-3bsuu01r
```

