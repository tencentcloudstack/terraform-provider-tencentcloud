Use this data source to query detailed information of teo zone configuration, including cache settings, origin configuration, security policies, and other zone-related configurations.

Example Usage

```hcl
data "tencentcloud_teo_export_zone_config" "export_config" {
  zone_id = "zone-2xkazzl8yf6k"
}

output "zone_name" {
  value = data.tencentcloud_teo_export_zone_config.export_config.zone_name
}

output "area" {
  value = data.tencentcloud_teo_export_zone_config.export_config.area
}

output "cache_config" {
  value = data.tencentcloud_teo_export_zone_config.export_config.cache_config
}

output "https_config" {
  value = data.tencentcloud_teo_export_zone_config.export_config.https
}
```

Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String, ForceNew) Zone ID.

Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `zone_name` - Zone name.
* `area` - Acceleration area. Values: `mainland`: Chinese mainland; `overseas`: Outside the Chinese mainland.
* `cache_key` - Node cache key configuration.
    * `switch` - Cache key configuration switch. Values: `on`: On; `off`: Off.
    * `ignore_case` - Whether to ignore case. Values: `on`: On; `off`: Off.
* `quic` - QUIC access configuration.
    * `switch` - QUIC configuration switch. Values: `on`: On; `off`: Off.
* `post_max_size` - POST request transmission configuration.
    * `switch` - POST request size limit configuration switch. Values: `on`: On; `off`: Off.
    * `max_size` - Maximum POST request size, unit: bytes.
* `compression` - Smart compression configuration.
    * `switch` - Compression configuration switch. Values: `on`: On; `off`: Off.
    * `algorithms` - Compression algorithm.
* `upstream_http2` - HTTP2 origin configuration.
    * `switch` - HTTP2 origin configuration switch. Values: `on`: On; `off`: Off.
* `force_redirect` - Access protocol forced HTTPS redirect configuration.
    * `switch` - Forced HTTPS redirect configuration switch. Values: `on`: On; `off`: Off.
    * `redirect_status_code` - Redirect status code.
* `cache_config` - Cache expiration time configuration.
    * `switch` - Cache configuration switch. Values: `on`: On; `off`: Off.
    * `rules` - Cache rule configuration list.
        * `cache_type` - Cache rule type.
        * `cache_time` - Cache time, unit: seconds.
* `origin` - Origin configuration.
    * `origins` - Origin list.
        * `origin_type` - Origin type.
        * `origin_value` - Origin address.
    * `origin_pull_protocol` - Origin pull protocol.
* `smart_routing` - Smart acceleration configuration.
    * `switch` - Smart acceleration configuration switch. Values: `on`: On; `off`: Off.
* `max_age` - Browser cache configuration.
    * `switch` - Browser cache configuration switch. Values: `on`: On; `off`: Off.
    * `follow_origin` - Whether to follow the origin server's Cache-Control header. Values: `on`: On; `off`: Off.
    * `max_age_time` - Browser cache time, unit: seconds.
* `offline_cache` - Offline cache configuration.
    * `switch` - Offline cache configuration switch. Values: `on`: On; `off`: Off.
* `websocket` - WebSocket configuration.
    * `switch` - WebSocket configuration switch. Values: `on`: On; `off`: Off.
    * `timeout` - WebSocket timeout time, unit: seconds.
* `client_ip_header` - Client IP origin request header configuration.
    * `switch` - Client IP header configuration switch. Values: `on`: On; `off`: Off.
    * `header_name` - Header name.
* `cache_prefresh` - Cache pre-refresh configuration.
    * `switch` - Cache pre-refresh configuration switch. Values: `on`: On; `off`: Off.
    * `percent` - Cache pre-refresh percentage.
* `ipv6` - IPv6 access configuration.
    * `switch` - IPv6 configuration switch. Values: `on`: On; `off`: Off.
* `https` - HTTPS acceleration configuration.
    * `switch` - HTTPS configuration switch. Values: `on`: On; `off`: Off.
    * `http2` - HTTP2 configuration.
        * `switch` - HTTP2 configuration switch. Values: `on`: On; `off`: Off.
    * `ocsp` - OCSP configuration.
        * `switch` - OCSP configuration switch. Values: `on`: On; `off`: Off.
    * `tls_version` - TLS version configuration.
        * `min_version` - Minimum TLS version.
        * `max_version` - Maximum TLS version.
        * `support_versions` - Supported TLS versions.
    * `cipher` - Cipher suite configuration.
        * `switch` - Custom cipher suite configuration switch. Values: `on`: On; `off`: Off.
        * `suites` - Cipher suite list.
* `client_ip_country` - Whether to carry client IP's region information when returning to origin.
    * `switch` - Client IP region information configuration switch. Values: `on`: On; `off`: Off.
    * `header_name` - Header name.
* `grpc` - GRPC protocol support configuration.
    * `switch` - GRPC configuration switch. Values: `on`: On; `off`: Off.
* `network_error_logging` - Network error log recording configuration.
    * `switch` - Network error log configuration switch. Values: `on`: On; `off`: Off.
* `image_optimize` - Image optimization related configuration.
    * `switch` - Image optimization configuration switch. Values: `on`: On; `off`: Off.
* `accelerate_mainland` - Chinese mainland acceleration optimization configuration.
    * `switch` - Chinese mainland acceleration optimization configuration switch. Values: `on`: On; `off`: Off.
* `standard_debug` - Standard Debug configuration.
    * `switch` - Standard Debug configuration switch. Values: `on`: On; `off`: Off.
* `jit_video_process` - Video instant processing configuration.
    * `switch` - Video instant processing configuration switch. Values: `on`: On; `off`: Off.
