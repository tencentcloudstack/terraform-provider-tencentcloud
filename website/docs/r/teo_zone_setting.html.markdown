---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_zone_setting"
sidebar_current: "docs-tencentcloud-resource-teo_zone_setting"
description: |-
  Provides a resource to create a teo zone_setting
---

# tencentcloud_teo_zone_setting

Provides a resource to create a teo zone_setting

## Example Usage

```hcl
resource "tencentcloud_teo_zone_setting" "zone_setting" {
  zone_id = "zone-297z8rf93cfw"

  cache {

    follow_origin {
      switch = "on"
    }

    no_cache {
      switch = "off"
    }
  }

  cache_key {
    full_url_cache = "on"
    ignore_case    = "off"

    query_string {
      action = "includeCustom"
      switch = "off"
      value  = []
    }
  }

  cache_prefresh {
    percent = 90
    switch  = "off"
  }

  client_ip_header {
    switch = "off"
  }

  compression {
    algorithms = [
      "brotli",
      "gzip",
    ]
    switch = "on"
  }

  force_redirect {
    redirect_status_code = 302
    switch               = "off"
  }

  https {
    http2         = "on"
    ocsp_stapling = "off"
    tls_version = [
      "TLSv1",
      "TLSv1.1",
      "TLSv1.2",
      "TLSv1.3",
    ]

    hsts {
      include_sub_domains = "off"
      max_age             = 0
      preload             = "off"
      switch              = "off"
    }
  }

  ipv6 {
    switch = "off"
  }

  max_age {
    follow_origin = "on"
    max_age_time  = 0
  }

  offline_cache {
    switch = "on"
  }

  origin {
    backup_origins       = []
    origin_pull_protocol = "follow"
    origins              = []
  }

  post_max_size {
    max_size = 524288000
    switch   = "on"
  }

  quic {
    switch = "off"
  }

  smart_routing {
    switch = "off"
  }

  upstream_http2 {
    switch = "off"
  }

  web_socket {
    switch  = "off"
    timeout = 30
  }
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String, ForceNew) Site ID.
* `cache_key` - (Optional, List) Node cache key configuration.
* `cache_prefresh` - (Optional, List) Cache pre-refresh configuration.
* `cache` - (Optional, List) Cache expiration time configuration.
* `client_ip_header` - (Optional, List) Origin-pull client IP header configuration.
* `compression` - (Optional, List) Smart compression configuration.
* `force_redirect` - (Optional, List) Force HTTPS redirect configuration.
* `https` - (Optional, List) HTTPS acceleration configuration.
* `ipv6` - (Optional, List) IPv6 access configuration.
* `max_age` - (Optional, List) Browser cache configuration.
* `offline_cache` - (Optional, List) Offline cache configuration.
* `origin` - (Optional, List) Origin server configuration.
* `post_max_size` - (Optional, List) Maximum size of files transferred over POST request.
* `quic` - (Optional, List) QUIC access configuration.
* `smart_routing` - (Optional, List) Smart acceleration configuration.
* `upstream_http2` - (Optional, List) HTTP2 origin-pull configuration.
* `web_socket` - (Optional, List) WebSocket configuration.

The `cache_key` object supports the following:

* `full_url_cache` - (Optional, String) Specifies whether to enable full-path cache.- `on`: Enable full-path cache (i.e., disable Ignore Query String).- `off`: Disable full-path cache (i.e., enable Ignore Query String). Note: This field may return null, indicating that no valid value can be obtained.
* `ignore_case` - (Optional, String) Specifies whether the cache key is case-sensitive. Note: This field may return null, indicating that no valid value can be obtained.
* `query_string` - (Optional, List) Request parameter contained in CacheKey. Note: This field may return null, indicating that no valid value can be obtained.

The `cache_prefresh` object supports the following:

* `switch` - (Required, String) Specifies whether to enable cache prefresh.- `on`: Enable.- `off`: Disable.
* `percent` - (Optional, Int) Percentage of cache time before try to prefresh. Valid value range: 1-99.

The `cache` object supports the following:

* `cache_time` - (Optional, Int) Cache expiration time settings.Unit: second. The maximum value is 365 days. Note: This field may return null, indicating that no valid value can be obtained.
* `ignore_cache_control` - (Optional, String) Specifies whether to enable force cache.- `on`: Enable.- `off`: Disable. Note: This field may return null, indicating that no valid value can be obtained.
* `switch` - (Optional, String) Cache configuration switch.- `on`: Enable.- `off`: Disable. Note: This field may return null, indicating that no valid value can be obtained.

The `cache` object supports the following:

* `cache` - (Optional, List) Cache configuration. Note: This field may return null, indicating that no valid value can be obtained.
* `follow_origin` - (Optional, List) Follows the origin server configuration. Note: This field may return null, indicating that no valid value can be obtained.
* `no_cache` - (Optional, List) No-cache configuration. Note: This field may return null, indicating that no valid value can be obtained.

The `client_ip_header` object supports the following:

* `switch` - (Required, String) Specifies whether to enable client IP header.- `on`: Enable.- `off`: Disable. Note: This field may return null, indicating that no valid value can be obtained.
* `header_name` - (Optional, String) Name of the origin-pull client IP request header. Note: This field may return null, indicating that no valid value can be obtained.

The `compression` object supports the following:

* `switch` - (Required, String) Whether to enable Smart compression.- `on`: Enable.- `off`: Disable.
* `algorithms` - (Optional, Set) Compression algorithms to select. Valid values: `brotli`, `gzip`.

The `follow_origin` object supports the following:

* `switch` - (Optional, String) Specifies whether to follow the origin server configuration.- `on`: Enable.- `off`: Disable. Note: This field may return null, indicating that no valid value can be obtained.

The `force_redirect` object supports the following:

* `switch` - (Required, String) Whether to enable force redirect.- `on`: Enable.- `off`: Disable.
* `redirect_status_code` - (Optional, Int) Redirection status code.- 301- 302 Note: This field may return null, indicating that no valid value can be obtained.

The `hsts` object supports the following:

* `switch` - (Required, String) - `on`: Enable.- `off`: Disable.
* `include_sub_domains` - (Optional, String) Specifies whether to include subdomain names. Valid values: `on` and `off`. Note: This field may return null, indicating that no valid value can be obtained.
* `max_age` - (Optional, Int) MaxAge value in seconds, should be no more than 1 day. Note: This field may return null, indicating that no valid value can be obtained.
* `preload` - (Optional, String) Specifies whether to preload. Valid values: `on` and `off`. Note: This field may return null, indicating that no valid value can be obtained.

The `https` object supports the following:

* `hsts` - (Optional, List) HSTS Configuration. Note: This field may return null, indicating that no valid value can be obtained.
* `http2` - (Optional, String) HTTP2 configuration switch.- `on`: Enable.- `off`: Disable. Note: This field may return null, indicating that no valid value can be obtained.
* `ocsp_stapling` - (Optional, String) OCSP configuration switch.- `on`: Enable.- `off`: Disable.It is disabled by default. Note: This field may return null, indicating that no valid value can be obtained.
* `tls_version` - (Optional, Set) TLS version settings. Valid values: `TLSv1`, `TLSV1.1`, `TLSV1.2`, and `TLSv1.3`.Only consecutive versions can be enabled at the same time. Note: This field may return null, indicating that no valid value can be obtained.

The `ipv6` object supports the following:

* `switch` - (Required, String) - `on`: Enable.- `off`: Disable.

The `max_age` object supports the following:

* `follow_origin` - (Optional, String) Specifies whether to follow the max cache age of the origin server.- `on`: Enable.- `off`: Disable.If it&#39;s on, MaxAgeTime is ignored. Note: This field may return null, indicating that no valid value can be obtained.
* `max_age_time` - (Optional, Int) Specifies the max age of the cache (in seconds). The maximum value is 365 days. Note: the value 0 means not to cache. Note: This field may return null, indicating that no valid value can be obtained.

The `no_cache` object supports the following:

* `switch` - (Optional, String) Whether to cache the configuration.- `on`: Do not cache.- `off`: Cache. Note: This field may return null, indicating that no valid value can be obtained.

The `offline_cache` object supports the following:

* `switch` - (Required, String) Whether to enable offline cache.- `on`: Enable.- `off`: Disable. Note: This field may return null, indicating that no valid value can be obtained.

The `origin` object supports the following:

* `backup_origins` - (Optional, Set) Backup origin sites list. Note: This field may return null, indicating that no valid value can be obtained.
* `cos_private_access` - (Optional, String) Whether access private cos bucket is allowed when `OriginType` is cos. Note: This field may return null, indicating that no valid value can be obtained.
* `origin_pull_protocol` - (Optional, String) Origin-pull protocol.- `http`: Switch HTTPS requests to HTTP.- `follow`: Follow the protocol of the request.- `https`: Switch HTTP requests to HTTPS. This only supports port 443 on the origin server. Note: This field may return null, indicating that no valid value can be obtained.
* `origins` - (Optional, Set) Origin sites list. Note: This field may return null, indicating that no valid value can be obtained.

The `post_max_size` object supports the following:

* `switch` - (Required, String) Specifies whether to enable custom setting of the maximum file size.- `on`: Enable. You can set a custom max size.- `off`: Disable. In this case, the max size defaults to 32 MB.
* `max_size` - (Optional, Int) Maximum size. Value range: 1-500 MB. Note: This field may return null, indicating that no valid value can be obtained.

The `query_string` object supports the following:

* `switch` - (Required, String) Whether to use QueryString as part of CacheKey.- `on`: Enable.- `off`: Disable. Note: This field may return null, indicating that no valid value can be obtained.
* `action` - (Optional, String) - `includeCustom`: Include the specified query strings.- `excludeCustom`: Exclude the specified query strings. Note: This field may return null, indicating that no valid value can be obtained.
* `value` - (Optional, Set) Array of query strings used/excluded. Note: This field may return null, indicating that no valid value can be obtained.

The `quic` object supports the following:

* `switch` - (Required, String) Whether to enable QUIC.- `on`: Enable.- `off`: Disable.

The `smart_routing` object supports the following:

* `switch` - (Required, String) Whether to enable smart acceleration.- `on`: Enable.- `off`: Disable.

The `upstream_http2` object supports the following:

* `switch` - (Required, String) Whether to enable HTTP2 origin-pull.- `on`: Enable.- `off`: Disable.

The `web_socket` object supports the following:

* `switch` - (Required, String) Whether to enable custom WebSocket timeout setting. When it&#39;s off: it means to keep the default WebSocket connection timeout period, which is 15 seconds. To change the timeout period, please set it to on.
* `timeout` - (Optional, Int) Sets timeout period in seconds. Maximum value: 120.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `area` - Acceleration area of the zone. Valid values: `mainland`, `overseas`.


## Import

teo zone_setting can be imported using the zone_id, e.g.
```
$ terraform import tencentcloud_teo_zone_setting.zone_setting zone-297z8rf93cfw#
```

