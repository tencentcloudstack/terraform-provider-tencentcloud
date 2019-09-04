---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_http_domain"
sidebar_current: "docs-tencentcloud-resource-gaap_http_domain"
description: |-
  Provides a resource to create a forward domain of layer7 listener.
---

# tencentcloud_gaap_http_domain

Provides a resource to create a forward domain of layer7 listener.

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
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
}

resource "tencentcloud_gaap_http_domain" "foo" {
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain      = "www.qq.com"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, ForceNew) Forward domain of the layer7 listener.
* `listener_id` - (Required, ForceNew) ID of the layer7 listener.
* `basic_auth_id` - (Optional) ID of the basic authentication.
* `basic_auth` - (Optional) Indicates whether basic authentication is enable, default is `false`.
* `certificate_id` - (Optional) ID of the server certificate, default value is `default`.
* `client_certificate_id` - (Optional) ID of the client certificate, default value is `default`.
* `gaap_auth_id` - (Optional) ID of the SSL certificate.
* `gaap_auth` - (Optional) Indicates whether SSL certificate authentication is enable, default is `false`.
* `realserver_auth` - (Optional) Indicates whether realserver authentication is enable, default is `false`.
* `realserver_certificate_domain` - (Optional) CA certificate domain of the realserver.
* `realserver_certificate_id` - (Optional) CA certificate ID of the realserver.


