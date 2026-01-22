Provides a resource to create a TEO l7 acc rule v2

Example Usage

```hcl
resource "tencentcloud_teo_l7_acc_rule_v2" "example" {
  zone_id     = "zone-3fkff38fyw8s"
  description = ["description"]
  rule_name   = "网站加速"
  status      = "enable"
  branches {
    condition = "$${http.request.host} in ['www.example.com']"
    actions {
      name = "Cache"
      cache_parameters {
        custom_time {
          cache_time           = 2592000
          ignore_cache_control = "off"
          switch               = "on"
        }
      }
    }

    actions {
      name = "CacheKey"
      cache_key_parameters {
        full_url_cache = "on"
        ignore_case    = "off"
        query_string {
          switch = "off"
          values = []
        }
      }
    }

    actions {
      name = "ModifyRequestHeader"
      modify_request_header_parameters {
        header_actions {
          action = "set"
          name   = "EO-Client-OS"
          value  = "*"
        }

        header_actions {
          action = "add"
          name   = "O-Client-Browser"
          value  = "*"
        }

        header_actions {
          action = "del"
          name   = "Eo-Client-Device"
        }
      }
    }

    actions {
      name = "ContentCompression"
      content_compression_parameters {
        switch = "on"
      }
    }

    sub_rules {
      description = ["1-1"]
      branches {
        condition = "lower($${http.request.file_extension}) in ['php', 'jsp', 'asp', 'aspx']"
        actions {
          name = "Cache"
          cache_parameters {
            no_cache {
              switch = "on"
            }
          }
        }
      }
    }

    sub_rules {
      description = ["1-2"]
      branches {
        condition = "$${http.request.file_extension} in ['jpg', 'png', 'gif', 'bmp', 'svg', 'webp']"
        actions {
          name = "MaxAge"
          max_age_parameters {
            cache_time    = 3600
            follow_origin = "off"
          }
        }
      }
    }
  }
}
```

Import

TEO l7 acc rule v2 can be imported using the {zone_id}#{rule_id}, e.g.

````
terraform import tencentcloud_teo_l7_acc_rule_v2.example zone-3fkff38fyw8s#rule-3ft1xeuhlj1b
````
