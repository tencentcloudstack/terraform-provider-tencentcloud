---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_http_domains"
sidebar_current: "docs-tencentcloud-datasource-gaap_http_domains"
description: |-
  Use this data source to query forward domain of layer7 listeners.
---

# tencentcloud_gaap_http_domains

Use this data source to query forward domain of layer7 listeners.

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

data "tencentcloud_gaap_http_domains" "foo" {
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain      = "${tencentcloud_gaap_http_domain.foo.domain}"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required) Forward domain of the layer7 listener to be queried.
* `listener_id` - (Required) ID of the layer7 listener to be queried.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `domains` - An information list of forward domain of the layer7 listeners. Each element contains the following attributes:
  * `basic_auth_id` - ID of the basic authentication.
  * `basic_auth` - Indicates whether basic authentication is enable.
  * `certificate_id` - ID of the server certificate.
  * `client_certificate_id` - (**Deprecated**) It has been deprecated from version 1.26.0. Use `client_certificate_ids` instead. ID of the client certificate.
  * `client_certificate_ids` - ID list of the client certificate.
  * `domain` - Forward domain of the layer7 listener.
  * `gaap_auth_id` - ID of the SSL certificate.
  * `gaap_auth` - Indicates whether SSL certificate authentication is enable.
  * `realserver_auth` - Indicates whether realserver authentication is enable.
  * `realserver_certificate_domain` - CA certificate domain of the realserver.
  * `realserver_certificate_id` - (**Deprecated**) It has been deprecated from version 1.28.0. Use `realserver_certificate_ids` instead. CA certificate ID of the realserver.
  * `realserver_certificate_ids` - CA certificate ID list of the realserver.


