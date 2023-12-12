Provides a resource to create a CDN domain.

~> **NOTE:** To disable most of configuration with switch, just modify switch argument to off instead of remove the whole block

Example Usage

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

    force_redirect {
      switch               = "on"
      redirect_type        = "http"
      redirect_status_code = 302
    }
  }

  tags = {
    hello = "world"
  }
}
```

Example Usage of cdn uses cache and request headers

```hcl
resource "tencentcloud_cdn_domain" "foo" {
  domain         = "xxxx.com"
  service_type   = "web"
  area           = "mainland"
  # full_url_cache = true # Deprecated, use cache_key below.
  cache_key {
    full_url_cache = "on"
  }
  range_origin_switch = "off"

  rule_cache{
  	cache_time = 10000
  	no_cache_switch="on"
  	re_validate="on"
  }

  request_header{
  	switch = "on"

  	header_rules {
  		header_mode = "add"
  		header_name = "tf-header-name"
  		header_value = "tf-header-value"
  		rule_type = "all"
  		rule_paths = ["*"]
  	}
  }

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

    force_redirect {
      switch               = "on"
      redirect_type        = "http"
      redirect_status_code = 302
    }
  }

  tags = {
    hello = "world"
  }
}
```

Example Usage of COS bucket url as origin

```hcl
resource "tencentcloud_cos_bucket" "bucket" {
  # Bucket format should be [custom name]-[appid].
  bucket = "demo-bucket-1251234567"
  acl    = "private"
}

# Create cdn domain
resource "tencentcloud_cdn_domain" "cdn" {
  domain         = "abc.com"
  service_type   = "web"
  area           = "mainland"
  # full_url_cache = false # Deprecated
  cache_key {
    full_url_cache = "off"
  }

  origin {
    origin_type          = "cos"
    origin_list          = [tencentcloud_cos_bucket.bucket.cos_bucket_url]
    server_name          = tencentcloud_cos_bucket.bucket.cos_bucket_url
    origin_pull_protocol = "follow"
    cos_private_access   = "on"
  }

  https_config {
    https_switch         = "off"
    http2_switch         = "off"
    ocsp_stapling_switch = "off"
    spdy_switch          = "off"
    verify_client        = "off"
  }
}
```

Import

CDN domain can be imported using the id, e.g.

```
$ terraform import tencentcloud_cdn_domain.foo xxxx.com
```