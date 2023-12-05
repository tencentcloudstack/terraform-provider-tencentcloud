Provides a resource to create a teo zone_setting

Example Usage

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
    switch     = "on"
  }

  force_redirect {
    redirect_status_code = 302
    switch               = "off"
  }

  https {
    http2         = "on"
    ocsp_stapling = "off"
    tls_version   = [
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
Import

teo zone_setting can be imported using the zone_id, e.g.
```
terraform import tencentcloud_teo_zone_setting.zone_setting zone-297z8rf93cfw#
```