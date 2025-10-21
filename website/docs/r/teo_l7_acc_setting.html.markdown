---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_l7_acc_setting"
sidebar_current: "docs-tencentcloud-resource-teo_l7_acc_setting"
description: |-
  Provides a resource to create a teo l7_acc_setting
---

# tencentcloud_teo_l7_acc_setting

Provides a resource to create a teo l7_acc_setting

## Example Usage

```hcl
resource "tencentcloud_teo_l7_acc_setting" "teo_l7_acc_setting" {
  zone_id = "zone-36bjhygh1bxe"
  zone_config {
    accelerate_mainland {
      switch = "on"
    }
    cache {
      custom_time {
        cache_time = 2592000
        switch     = "off"
      }
      follow_origin {
        default_cache          = "off"
        default_cache_strategy = "on"
        default_cache_time     = 0
        switch                 = "on"
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
      }
    }
    cache_prefresh {
      cache_time_percent = 90
      switch             = "off"
    }
    client_ip_country {
      switch = "off"
    }
    client_ip_header {
      switch = "off"
    }
    compression {
      algorithms = ["brotli", "gzip"]
      switch     = "on"
    }
    force_redirect_https {
      redirect_status_code = 302
      switch               = "off"
    }
    grpc {
      switch = "off"
    }
    hsts {
      include_sub_domains = "off"
      preload             = "off"
      switch              = "off"
      timeout             = 0
    }
    http2 {
      switch = "off"
    }
    ipv6 {
      switch = "off"
    }
    max_age {
      cache_time    = 600
      follow_origin = "on"
    }
    ocsp_stapling {
      switch = "off"
    }
    offline_cache {
      switch = "on"
    }
    post_max_size {
      max_size = 838860800
      switch   = "on"
    }
    quic {
      switch = "off"
    }
    smart_routing {
      switch = "off"
    }
    standard_debug {
      allow_client_ip_list = []
      expires              = "1969-12-31T16:00:00Z"
      switch               = "off"
    }
    tls_config {
      cipher_suite = "loose-v2023"
      version      = ["TLSv1", "TLSv1.1", "TLSv1.2", "TLSv1.3"]
    }
    upstream_http2 {
      switch = "off"
    }
    web_socket {
      switch  = "off"
      timeout = 30
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `zone_config` - (Required, List) Site acceleration global configuration. the settings in this parameter will apply to all domain names under the site. you only need to modify the required settings directly, and other settings not passed in will remain unchanged.
* `zone_id` - (Required, String, ForceNew) Zone id.

The `accelerate_mainland` object of `zone_config` supports the following:

* `switch` - (Optional, String) Mainland china acceleration optimization switch. valid values:
on: Enable;
off: Disable.

The `cache_key` object of `zone_config` supports the following:

* `full_url_cache` - (Optional, String) Whether to enable full-path cache. values:
on: Enable full-path cache (i.e., disable ignore query string);
off: Disable full-path cache (i.e., enable ignore query string).
* `ignore_case` - (Optional, String) Whether to ignore case in the cache key. values:
on: Ignore;
off: Not ignore.
* `query_string` - (Optional, List) Query string retention configuration parameter. this field and fullurlcache must be set simultaneously, but cannot both be on.

The `cache_prefresh` object of `zone_config` supports the following:

* `cache_time_percent` - (Optional, Int) Prefresh interval set as a percentage of the node cache time. value range: 1-99.
Note: This field is required when switch is on; when switch is off, this field is not required and will not take effect if filled.
* `switch` - (Optional, String) Whether to enable cache prefresh. values:
on: Enable;
off: Disable.

The `cache` object of `zone_config` supports the following:

* `custom_time` - (Optional, List) Custom cache time configuration. only one of followorigin, nocache, customtime can have switch set to on.
* `follow_origin` - (Optional, List) Follow origin server cache configuration. only one of followorigin, nocache, customtime can have switch set to on.
* `no_cache` - (Optional, List) No cache configuration. only one of followorigin, nocache, customtime can have switch set to on.

The `client_ip_country` object of `zone_config` supports the following:

* `header_name` - (Optional, String) Name of the request header that contains the client IP region. It is valid when Switch=on.
The default value EO-Client-IPCountry is used when it is not specified.
* `switch` - (Optional, String) Whether to enable configuration. Values:
on: Enable;
off: Disable.

The `client_ip_header` object of `zone_config` supports the following:

* `header_name` - (Optional, String) Name of the request header containing the client ip address for origin-pull. when switch is on, this parameter is required. x-forwarded-for is not allowed for this parameter.
* `switch` - (Optional, String) Whether to enable configuration. values:
on: Enable;
off: Disable.

The `compression` object of `zone_config` supports the following:

* `algorithms` - (Optional, Set) Supported compression algorithm list. this field is required when switch is on; otherwise, it is not effective. valid values:
brotli: Brotli algorithm;
gzip: Gzip algorithm.
* `switch` - (Optional, String) Whether to enable smart compression. values:
on: Enable;
off: Disable.

The `custom_time` object of `cache` supports the following:

* `cache_time` - (Optional, Int) Custom cache time value, unit: seconds. value range: 0-315360000.
Note: This field is required when switch is on; when switch is off, this field is not required and will not take effect if filled.
* `switch` - (Optional, String) Custom cache time switch. values:
on: Enable;
off: Disable.

The `follow_origin` object of `cache` supports the following:

* `switch` - (Required, String) Whether to enable the configuration of following the origin server. Valid values:
on: Enable;
off: Disable.
* `default_cache_strategy` - (Optional, String) Whether to use the default caching policy when an origin server does not return the cache-control header. this field is required when defaultcache is set to on; otherwise, it is ineffective. when defaultcachetime is not 0, this field should be off. valid values:
on: Use the default caching policy.
off: Do not use the default caching policy.
* `default_cache_time` - (Optional, Int) The default cache time in seconds when an origin server does not return the cache-control header. the value ranges from 0 to 315360000. this field is required when defaultcache is set to on; otherwise, it is ineffective. when defaultcachestrategy is on, this field should be 0.
* `default_cache` - (Optional, String) Whether to cache when an origin server does not return the cache-control header. this field is required when switch is on; when switch is off, this field is not required and will be ineffective if filled. valid values:
on: Cache;
off: Do not cache.

The `force_redirect_https` object of `zone_config` supports the following:

* `redirect_status_code` - (Optional, Int) Redirection status code. this field is required when switch is on; otherwise, it is not effective. valid values are:
301: 301 redirect;
302: 302 redirect.
* `switch` - (Optional, String) Whether to enable forced redirect configuration switch. values:
on: Enable;
off: Disable.

The `grpc` object of `zone_config` supports the following:

* `switch` - (Optional, String) Whether to enable grpc. values:
on: Enable;
off: Disable.

The `hsts` object of `zone_config` supports the following:

* `include_sub_domains` - (Optional, String) Whether to allow other subdomains to inherit the same hsts header. values:
on: Allows other subdomains to inherit the same hsts header;
off: Does not allow other subdomains to inherit the same hsts header.
Note: When switch is on, this field is required; when switch is off, this field is not required and will not take effect if filled.
* `preload` - (Optional, String) Whether to allow the browser to preload the hsts header. valid values:
on: Allows the browser to preload the hsts header;
off: Does not allow the browser to preload the hsts header.
Note: When switch is on, this field is required; when switch is off, this field is not required and will not take effect if filled.
* `switch` - (Optional, String) Whether to enable hsts. values:
on: Enable;
off: Disable.
* `timeout` - (Optional, Int) Cache hsts header time, unit: seconds. value range: 1-31536000.
Note: This field is required when switch is on; when switch is off, this field is not required and will not take effect if filled.

The `http2` object of `zone_config` supports the following:

* `switch` - (Optional, String) Whether to enable http2 access. values:
on: Enable;
off: Disable.

The `ipv6` object of `zone_config` supports the following:

* `switch` - (Optional, String) Whether to enable ipv6 access functionality. valid values:
on: Enable;
off: Disable.

The `max_age` object of `zone_config` supports the following:

* `cache_time` - (Optional, Int) Custom cache time value, unit: seconds. value range: 0-315360000.
Note: When followorigin is off, it means not following the origin server and using cachetime to set the cache time; otherwise, this field will not take effect.
* `follow_origin` - (Optional, String) Specifies whether to follow the origin server cache-control configuration, with the following values:
on: Follow the origin server and ignore the field cachetime;
off: Do not follow the origin server and apply the field cachetime.

The `no_cache` object of `cache` supports the following:

* `switch` - (Required, String) Whether to enable no-cache configuration. Valid values:
on: Enable;
off: Disable.

The `ocsp_stapling` object of `zone_config` supports the following:

* `switch` - (Optional, String) Whether to enable ocsp stapling configuration switch. values:
on: Enable;
off: Disable.

The `offline_cache` object of `zone_config` supports the following:

* `switch` - (Optional, String) Whether to enable offline caching. values:
on: Enable;
off: Disable.

The `post_max_size` object of `zone_config` supports the following:

* `max_size` - (Optional, Int) Maximum size of the file uploaded for streaming via a post request, in bytes. value range: 1 * 2^20 bytes to 500 * 2^20 bytes.
* `switch` - (Optional, String) Whether to enable post request file upload limit, in bytes (default limit: 32 * 2^20 bytes). valid values:
on: Enable limit;
off: Disable limit.

The `query_string` object of `cache_key` supports the following:

* `action` - (Optional, String) Actions to retain/ignore specified parameters in the query string. values:
includeCustom: retain partial parameters;
excludeCustom: ignore partial parameters.
Note: This field is required when switch is on. when switch is off, this field is not required and will not take effect if filled.
* `switch` - (Optional, String) Query string retain/ignore specified parameter switch. valid values are:
on: Enable;
off: Disable.
* `values` - (Optional, Set) List of parameter names to be retained/ignored in the query string.
note: This field is required when switch is on. when switch is off, this field is not required and will not take effect if filled.

The `quic` object of `zone_config` supports the following:

* `switch` - (Optional, String) Whether to enable quic. values:
on: Enable;
off: Disable.

The `smart_routing` object of `zone_config` supports the following:

* `switch` - (Optional, String) Whether to enable smart acceleration. values:
on: Enable;
off: Disable.

The `standard_debug` object of `zone_config` supports the following:

* `allow_client_ip_list` - (Optional, Set) The client ip to allow. it can be an ipv4/ipv6 address or a cidr block. `0.0.0.0/0` means to allow all ipv4 clients for debugging; `::/0` means to allow all ipv6 clients for debugging; `127.0.0.1` is not allowed.
Note: this field is required when switch=on and the number of entries should be 1-100. when switch=off, this field is not required and any value specified will not take effect.
* `expires` - (Optional, String) Debug feature expiration time. the feature will be disabled after the set time.
Note: this field is required when switch=on. when switch=off, this field is not required and any value specified will not take effect.
* `switch` - (Optional, String) Whether to enable standard debugging. values:
on: Enable;
off: Disable.

The `tls_config` object of `zone_config` supports the following:

* `cipher_suite` - (Optional, String) Cipher suite. for detailed information, please refer to tls versions and cipher suites description. valid values:
loose-v2023: loose-v2023 cipher suite;
general-v2023: general-v2023 cipher suite;
strict-v2023: strict-v2023 cipher suite.
* `version` - (Optional, Set) TLS version. at least one must be specified. if multiple versions are specified, they must be consecutive, e.g., enable tls1, 1.1, 1.2, and 1.3. it is not allowed to enable only 1 and 1.2 while disabling 1.1. valid values:
TLSv1: TLSv1 version;
TLSv1.1: TLSv1.1 version;
TLSv1.2: TLSv1.2 version;
TLSv1.3: TLSv1.3 version.

The `upstream_http2` object of `zone_config` supports the following:

* `switch` - (Optional, String) Whether to enable http2 origin-pull. valid values:
on: Enable;
off: Disable.

The `web_socket` object of `zone_config` supports the following:

* `switch` - (Optional, String) Whether to enable websocket connection timeout. values:
on: Use timeout as the websocket timeout;
off: The platform still supports websocket connections, using the system default timeout of 15 seconds.
* `timeout` - (Optional, Int) Timeout, unit: seconds. maximum timeout is 120 seconds.
Note: This field is required when switch is on; otherwise, this field will not take effect.

The `zone_config` object supports the following:

* `accelerate_mainland` - (Optional, List) Accelerate optimization and configuration in mainland china.
* `cache_key` - (Optional, List) The node cache key configuration.
* `cache_prefresh` - (Optional, List) Cache prefresh configuration.
* `cache` - (Optional, List) Node cache expiration time configuration.
* `client_ip_country` - (Optional, List) Client ip origin-pull request header configuration.
* `client_ip_header` - (Optional, List) Client ip origin-pull request header configuration.
* `compression` - (Optional, List) Smart compression configuration.
* `force_redirect_https` - (Optional, List) Forced https redirect configuration for access protocols.
* `grpc` - (Optional, List) Configuration of grpc support.
* `hsts` - (Optional, List) HSTS configuration.
* `http2` - (Optional, List) HTTP/2 configuration.
* `ipv6` - (Optional, List) IPv6 access configuration.
* `max_age` - (Optional, List) Browser cache rule configuration, which is used to set the default value of maxage and is disabled by default.
* `ocsp_stapling` - (Optional, List) OCSP stapling configuration.
* `offline_cache` - (Optional, List) Offline cache configuration.
* `post_max_size` - (Optional, List) POST request transport configuration.
* `quic` - (Optional, List) QUIC access configuration.
* `smart_routing` - (Optional, List) Intelligent acceleration configuration.
* `standard_debug` - (Optional, List) Standard debugging configuration.
* `tls_config` - (Optional, List) TLS configuration.
* `upstream_http2` - (Optional, List) HTTP2 origin-pull configuration.
* `web_socket` - (Optional, List) WebSocket configuration.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `zone_name` - Zone name.


## Import

teo l7_acc_setting can be imported using the zone_id, e.g.
````
terraform import tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting zone-297z8rf93cfw
````

