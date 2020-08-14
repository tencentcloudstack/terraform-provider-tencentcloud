---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdn_domains"
sidebar_current: "docs-tencentcloud-datasource-cdn_domains"
description: |-
  Use this data source to query the detail information of CDN domain.
---

# tencentcloud_cdn_domains

Use this data source to query the detail information of CDN domain.

## Example Usage

```hcl
data "tencentcloud_cdn_domains" "foo" {
  domain               = "xxxx.com"
  service_type         = "web"
  full_url_cache       = false
  origin_pull_protocol = "follow"
  https_switch         = "on"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Optional) Acceleration domain name.
* `full_url_cache` - (Optional) Whether to enable full-path cache.
* `https_switch` - (Optional) HTTPS configuration. The available value include `on`, `off` and `processing`.
* `origin_pull_protocol` - (Optional) Origin-pull protocol configuration. The available value include `http`, `https` and `follow`.
* `result_output_file` - (Optional) Used to save results.
* `service_type` - (Optional) Service type of acceleration domain name. The available value include `web`, `download` and `media`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `domain_list` - An information list of cdn domain. Each element contains the following attributes:
  * `area` - Acceleration region.
  * `cname` - CNAME address of domain name.
  * `create_time` - Domain name creation time.
  * `domain` - Acceleration domain name.
  * `full_url_cache` - Whether to enable full-path cache.
  * `https_config` - HTTPS acceleration configuration. It's a list and consist of at most one item.
    * `http2_switch` - HTTP2 configuration switch.
    * `https_switch` - HTTPS configuration switch.
    * `ocsp_stapling_switch` - OCSP configuration switch.
    * `spdy_switch` - Spdy configuration switch.
    * `verify_client` - Client certificate authentication feature.
  * `id` - Domain name ID.
  * `origin` - Origin server configuration.
    * `backup_origin_list` - Backup origin server list.
    * `backup_origin_type` - Backup origin server type.
    * `backup_server_name` - Host header used when accessing the backup origin server. If left empty, the ServerName of master origin server will be used by default.
    * `cos_private_access` - When OriginType is COS, you can specify if access to private buckets is allowed.
    * `origin_list` - Master origin server list.
    * `origin_pull_protocol` - Origin-pull protocol configuration.
    * `origin_type` - Master origin server type.
    * `server_name` - Host header used when accessing the master origin server. If left empty, the acceleration domain name will be used by default.
  * `project_id` - The project CDN belongs to.
  * `service_type` - Service type of acceleration domain name.
  * `status` - Acceleration service status.
  * `tags` - Tags of cdn domain.
  * `update_time` - Last modified time of domain name.


