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


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateDomain](https://cloud.tencent.com/document/api/608/36991) - This API (CreateDomain) is used to create the access domain name for the HTTP/HTTPS listener. Clients request the backend data by accessing this domain.
* Create: [DescribeHTTPListeners](https://cloud.tencent.com/document/api/608/37002) - This API (DescribeHTTPListeners) is used to query HTTP listener information.
* Create: [DescribeHTTPSListeners](https://cloud.tencent.com/document/api/608/37001) - This API (DescribeHTTPSListeners) is used to query HTTPS listener information.
* Create: [DescribeRules](https://cloud.tencent.com/document/api/608/36986) - This API (DescribeRules) is used to query all rule information of a listener, including the domain names, paths, and list of bound origin servers. For version 3.0 connections, this API returns the advanced authentication configuration information of the domain name.
* Create: [SetAuthentication](https://cloud.tencent.com/document/api/608/36982) - This API (SetAuthentication) is used for the advanced authentication configuration of connections, including authentication methods and their certificates. If only supports connections of version 3.0.
* Read: [DescribeRules](https://cloud.tencent.com/document/api/608/36986) - This API (DescribeRules) is used to query all rule information of a listener, including the domain names, paths, and list of bound origin servers. For version 3.0 connections, this API returns the advanced authentication configuration information of the domain name.
* Update: [DescribeHTTPSListeners](https://cloud.tencent.com/document/api/608/37001) - This API (DescribeHTTPSListeners) is used to query HTTPS listener information.
* Update: [ModifyCertificate](https://cloud.tencent.com/document/api/608/36985) - This API (ModifyCertificate) is used to modify a domain name certificate of a listener. domain name certificate. This API is only applicable to connections of version 3.0.
* Update: [ModifyDomain](https://cloud.tencent.com/document/api/608/36984) - This API (ModifyDomain) is used to modify domain names of listeners. For connections of version 3.0, it supports modifying certificates of the domain names.
* Update: [SetAuthentication](https://cloud.tencent.com/document/api/608/36982) - This API (SetAuthentication) is used for the advanced authentication configuration of connections, including authentication methods and their certificates. If only supports connections of version 3.0.
* Delete: [DeleteDomain](https://cloud.tencent.com/document/api/608/36989) - This API (DeleteDomain) is only applicable to layer-7 listeners. It is used to delete the domain names of this listener, and all the rules under the domain name. All rules bound to the origin server are unbound automatically.
* Delete: [DescribeRules](https://cloud.tencent.com/document/api/608/36986) - This API (DescribeRules) is used to query all rule information of a listener, including the domain names, paths, and list of bound origin servers. For version 3.0 connections, this API returns the advanced authentication configuration information of the domain name.

## Import

GAAP http domain can be imported using the id, e.g.

-> **NOTE:** The format of tencentcloud_gaap_http_domain id is `[listener-id]+[protocol]+[domain]`.

```
  $ terraform import tencentcloud_gaap_http_domain.foo listener-11112222+HTTP+www.qq.com
```

