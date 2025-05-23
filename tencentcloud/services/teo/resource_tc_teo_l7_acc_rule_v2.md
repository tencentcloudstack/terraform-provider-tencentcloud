Provides a resource to create a TEO l7 acc rule

Example Usage

```hcl
resource "tencentcloud_teo_l7_acc_rule_v2" "teo_l7_acc_rule_v2" {
  zone_id     = "zone-39quuimqg8r6"
  description = ["1"]
  rule_name   = "网站加速1"
  status = "enable"
  branches {
    condition = "$${http.request.host} in ['aaa.makn.cn']"
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

TEO l7 acc rule can be imported using the id, e.g.

````
terraform import tencentcloud_teo_l7_acc_rule.example ${zone_id}#${rule_id}
````
