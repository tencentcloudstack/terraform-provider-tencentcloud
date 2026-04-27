Provides a resource to create a teo l7_acc_setting

Example Usage

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
      switch      = "off"
    }
    client_ip_header {
      switch      = "off"
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

Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) Zone id.
* `zone_config` - (Required) Site acceleration global configuration.

The `zone_config` object supports the following:

* `smart_routing` - (Optional) Intelligent acceleration configuration.
* `cache` - (Optional) Node cache expiration time configuration.
* `max_age` - (Optional) Browser cache rule configuration.
* `cache_key` - (Optional) The node cache key configuration.
* `cache_prefresh` - (Optional) Cache prefresh configuration.
* `offline_cache` - (Optional) Offline cache configuration.
* `compression` - (Optional) Smart compression configuration.
* `force_redirect_https` - (Optional) Forced https redirect configuration for access protocols.
* `hsts` - (Optional) HSTS configuration.
* `tls_config` - (Optional) TLS configuration.
* `ocsp_stapling` - (Optional) OCSP stapling configuration.
* `http2` - (Optional) HTTP/2 configuration.
* `quic` - (Optional) QUIC access configuration.
* `upstream_http2` - (Optional) HTTP2 origin-pull configuration.
* `ipv6` - (Optional) IPv6 access configuration.
* `web_socket` - (Optional) WebSocket configuration.
* `post_max_size` - (Optional) POST request transport configuration.
* `client_ip_header` - (Optional) Client ip origin-pull request header configuration.
* `client_ip_country` - (Optional) Client ip origin-pull request header configuration.
* `grpc` - (Optional) Configuration of grpc support.
* `accelerate_mainland` - (Optional) Accelerate optimization and configuration in mainland china.
* `standard_debug` - (Optional) Standard debugging configuration.

Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `zone_name` - Zone name.
* `zone_setting` - Site acceleration configuration computed from the API response.
  * `zone_name` - Zone name.
  * `zone_config` - Site acceleration global configuration.
    * `smart_routing` - Intelligent acceleration configuration.
    * `cache` - Node cache expiration time configuration.
    * `max_age` - Browser cache rule configuration.
    * `cache_key` - The node cache key configuration.
    * `cache_prefresh` - Cache prefresh configuration.
    * `offline_cache` - Offline cache configuration.
    * `compression` - Smart compression configuration.
    * `force_redirect_https` - Forced https redirect configuration for access protocols.
    * `hsts` - HSTS configuration.
    * `tls_config` - TLS configuration.
    * `ocsp_stapling` - OCSP stapling configuration.
    * `http2` - HTTP/2 configuration.
    * `quic` - QUIC access configuration.
    * `upstream_http2` - HTTP2 origin-pull configuration.
    * `ipv6` - IPv6 access configuration.
    * `web_socket` - WebSocket configuration.
    * `post_max_size` - POST request transport configuration.
    * `client_ip_header` - Client ip origin-pull request header configuration.
    * `client_ip_country` - Client ip origin-pull request header configuration.
    * `grpc` - Configuration of grpc support.
    * `accelerate_mainland` - Accelerate optimization and configuration in mainland china.
    * `standard_debug` - Standard debugging configuration.

Import

teo l7_acc_setting can be imported using the zone_id, e.g.
````
terraform import tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting zone-297z8rf93cfw
````