---
subcategory: "Global Application Acceleration(GAAP)"
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
  proxy_id = tencentcloud_gaap_proxy.foo.id
}

resource "tencentcloud_gaap_http_domain" "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Forward domain of the layer7 listener.
* `listener_id` - (Required, String, ForceNew) ID of the layer7 listener.
* `basic_auth_id` - (Optional, String) ID of the basic authentication.
* `basic_auth` - (Optional, Bool) Indicates whether basic authentication is enable, default value is `false`.
* `certificate_id` - (Optional, String) ID of the server certificate, default value is `default`.
* `client_certificate_id` - (Optional, String, **Deprecated**) It has been deprecated from version 1.26.0. Set `client_certificate_ids` instead. ID of the client certificate, default value is `default`.
* `client_certificate_ids` - (Optional, Set: [`String`]) ID list of the poly client certificate.
* `gaap_auth_id` - (Optional, String) ID of the SSL certificate.
* `gaap_auth` - (Optional, Bool) Indicates whether SSL certificate authentication is enable, default value is `false`.
* `group_id` - (Optional, String, ForceNew) Group Id.
* `is_default_server` - (Optional, Bool) Whether to use as the default domain name, the default is false.
* `realserver_auth` - (Optional, Bool) Indicates whether realserver authentication is enable, default value is `false`.
* `realserver_certificate_domain` - (Optional, String) CA certificate domain of the realserver. It has been deprecated.
* `realserver_certificate_id` - (Optional, String, **Deprecated**) It has been deprecated from version 1.28.0. Set `realserver_certificate_ids` instead. CA certificate ID of the realserver.
* `realserver_certificate_ids` - (Optional, Set: [`String`]) CA certificate ID list of the realserver.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

GAAP http domain can be imported using the id, e.g.

-> **NOTE:** The format of tencentcloud_gaap_http_domain id is `[listener-id]+[protocol]+[domain]`.

```
  $ terraform import tencentcloud_gaap_http_domain.foo listener-11112222+HTTP+www.qq.com
```

