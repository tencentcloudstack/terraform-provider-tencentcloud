resource "tencentcloud_teo_zone" "example" {
  zone_name = "example.com"
  plan_type = "<your-plan-type>"

  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_teo_zone_setting" "example" {
  zone_id = tencentcloud_teo_zone.example.id

  # Cache Configuration
  cache {
    follow_origin {
      switch = "on"
    }
  }
  # CacheKey Configuration
  cache_key {
    full_url_cache = "off"
    ignore_case    = "on"
    query_string {
      switch = "on"
      action = "includeCustom" # use specific parameters from URL
      value  = ["param0", "param1"]
    }
  }
  # HTTPS Configuration
  https {
    http2         = "on"
    ocsp_stapling = "on"
    tls_version   = ["TLSv1.2", "TLSv1.3"]
    hsts {
      include_sub_domains = "off"
      max_age             = 0
      preload             = "off"
      switch              = "off"
    }
  }
  # Smart Compression Configuration
  compression {
    switch     = "on"
    algorithms = ["brotli", "gzip"]
  }
  # Carry client IP to origin site
  client_ip_header {
    switch      = "on"
    header_name = "EO-Client-IPCountry"
  }
}

resource "tencentcloud_teo_dns_record" "rule_record" {
  zone_id = tencentcloud_teo_zone.example.id
  type    = "A"
  name    = "rule.example.com"
  content = "1.1.1.1"
  mode    = "proxied"
  ttl     = 300
}

# subdomain specific configuration
resource "tencentcloud_teo_rule_engine" "rule_example" {
  zone_id   = tencentcloud_teo_zone.example.id
  rule_name = "example_rule"
  status    = "enable"

  rules {
    # when request host is rule.example.com and file suffix is mp3 or mp4
    or {
      and {
        target   = "host"
        operator = "equal"
        values   = [tencentcloud_teo_dns_record.rule_record.name]
      }
      and {
        target   = "extension"
        operator = "equal"
        values   = ["mp4", "mp3"]
      }
    }

    actions {
      normal_action {
        action = "CacheKey"
        # CacheKey is ignore case
        parameters {
          name   = "Type"
          values = ["IgnoreCase"]
        }
        parameters {
          name   = "Switch"
          values = ["off"]
        }
        # CacheKey should use User-Agent Header
        parameters {
          name   = "Type"
          values = ["Header"]
        }
        parameters {
          name   = "Switch"
          values = ["on"]
        }
        parameters {
          name   = "Value"
          values = ["User-Agent"]
        }
      }
    }

    # Add a HTTP header to response
    actions {
      rewrite_action {
        action = "ResponseHeader"
        parameters {
          action = "add"
          name   = "Added-Header"
          values = ["Added-Value"]
        }
      }
    }
  }
}