---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdn_domain"
sidebar_current: "docs-tencentcloud-resource-cdn_domain"
description: |-
  Provides a resource to create a CDN domain.
---

# tencentcloud_cdn_domain

Provides a resource to create a CDN domain.

## Example Usage

```hcl
resource "tencentcloud_cdn_domain" "foo" {
  domain         = "xxxx.com"
  service_type   = "web"
  area           = "mainland"
  full_url_cache = false

  origin {
    origin_type          = "ip"
    origin_list          = ["127.0.0.1"]
    origin_pull_protocol = "follow"
  }

  https_config {
    https_switch         = "off"
    http2_switch         = "off"
    ocsp_stapling_switch = "off"
    spdy_switch          = "off"
    verify_client        = "off"
  }

  tags = {
    hello = "world"
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, ForceNew) Name of the acceleration domain.
* `origin` - (Required) Origin server configuration. It's a list and consist of at most one item.
* `service_type` - (Required, ForceNew) Service type of Acceleration domain name. Valid values are `web`, `download` and `media`.
* `area` - (Optional) Domain name acceleration region.  Valid values are `mainland`, `overseas` and `global`.
* `full_url_cache` - (Optional) Whether to enable full-path cache. Default value is `true`.
* `https_config` - (Optional) HTTPS acceleration configuration. It's a list and consist of at most one item.
* `project_id` - (Optional) The project CDN belongs to, default to 0.
* `tags` - (Optional) Tags of cdn domain.

The `client_certificate_config` object supports the following:

* `certificate_content` - (Required) Client Certificate PEM format, requires Base64 encoding.

The `https_config` object supports the following:

* `https_switch` - (Required) HTTPS configuration switch. Valid values are `on` and `off`.
* `client_certificate_config` - (Optional) Client certificate configuration information.
* `http2_switch` - (Optional) HTTP2 configuration switch. Valid values are `on` and `off`, and default value is `off`.
* `ocsp_stapling_switch` - (Optional) OCSP configuration switch. Valid values are `on` and `off`, and default value is `off`.
* `server_certificate_config` - (Optional) Server certificate configuration information.
* `spdy_switch` - (Optional) Spdy configuration switch. Valid values are `on` and `off`, and default value is `off`.
* `verify_client` - (Optional) Client certificate authentication feature. Valid values are `on` and `off`, and default value is `off`.

The `origin` object supports the following:

* `origin_list` - (Required) Master origin server list. Valid values can be ip or doamin name. When modifying the origin server, you need to enter the corresponding `origin_type`.
* `origin_type` - (Required) Master origin server type. Valid values are `domain`, `cos`, `ip`, `ipv6` and `ip_ipv6`.
* `backup_origin_list` - (Optional) Backup origin server list. Valid values can be ip or doamin name. When modifying the backup origin server, you need to enter the corresponding `backup_origin_type`.
* `backup_origin_type` - (Optional) Backup origin server type. Valid values are `domain` and `ip`.
* `backup_server_name` - (Optional) Host header used when accessing the backup origin server. If left empty, the ServerName of master origin server will be used by default.
* `cos_private_access` - (Optional) When OriginType is COS, you can specify if access to private buckets is allowed. Valid values are `on` and `off`, and default value is `off`.
* `origin_pull_protocol` - (Optional) Origin-pull protocol configuration. Valid values are `http`, `https` and `follow`, and default value is `http`.
* `server_name` - (Optional) Host header used when accessing the master origin server. If left empty, the acceleration domain name will be used by default.

The `server_certificate_config` object supports the following:

* `certificate_content` - (Optional) Server certificate information. This is required when uploading an external certificate, which should contain the complete certificate chain.
* `certificate_id` - (Optional) Server certificate ID.
* `message` - (Optional) Certificate remarks.
* `private_key` - (Optional) Server key information. This is required when uploading an external certificate.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cname` - CNAME address of domain name.
* `create_time` - Creation time of domain name.
* `status` - Acceleration service status.


## Import

CDN domain can be imported using the id, e.g.

```
$ terraform import tencentcloud_cdn_domain.foo xxxx.com
```

