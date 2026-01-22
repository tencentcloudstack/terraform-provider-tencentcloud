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
Import

teo l7_acc_setting can be imported using the zone_id, e.g.
````
terraform import tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting zone-297z8rf93cfw
````