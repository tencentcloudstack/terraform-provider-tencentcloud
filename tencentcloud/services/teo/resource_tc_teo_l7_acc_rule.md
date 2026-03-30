Provides a resource to create a TEO l7 acc rule

~> **NOTE:** This feature only supports the sites in the plans of the Standard Edition and the Enterprise Edition.

~> **NOTE:** The `filters` parameter can be used to filter query results by rule name, status, or other supported fields.

Example Usage

```hcl
resource "tencentcloud_teo_l7_acc_rule" "example" {
  zone_id = "zone-36bjhygh1bxe"
  rules {
    description = ["1"]
    rule_name   = "Web Acceleration"
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
          scheme         = null
          query_string {
            action = null
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

  rules {
    description = ["2"]
    rule_name   = "Live Video Streaming"
    branches {
      condition = "$${http.request.host} in ['aaa.makn.cn']"
      sub_rules {
        description = ["2-1"]
        branches {
          condition = "$${http.request.file_extension} in ['m3u8', 'mpd']"
          actions {
            name = "Cache"
            cache_parameters {
              custom_time {
                cache_time           = 1
                ignore_cache_control = "off"
                switch               = "on"
              }
            }
          }
        }

        branches {
          condition = "$${http.request.file_extension} in ['ts', 'mp4', 'm4a', 'm4s']"
          actions {
            name = "Cache"
            cache_parameters {
              custom_time {
                cache_time           = 86400
                ignore_cache_control = "off"
                switch               = "on"
              }
            }
          }
        }

        branches {
          condition = "*"
          actions {
            name = "Cache"
            cache_parameters {
              follow_origin {
                default_cache          = "on"
                default_cache_strategy = "on"
                default_cache_time     = 0
                switch                 = "on"
              }
            }
          }
        }
      }
    }
  }

  rules {
    description = ["3"]
    rule_name   = "Large File Download"
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
          full_url_cache = "off"
          ignore_case    = null
          scheme         = null
          query_string {
            action = null
            switch = "off"
            values = []
          }
        }
      }

      actions {
        name = "RangeOriginPull"
        range_origin_pull_parameters {
          switch = "on"
        }
      }

      sub_rules {
        description = ["3-1"]
        branches {
          condition = "$${http.request.file_extension} in ['php', 'jsp', 'asp', 'aspx']"
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
    }
  }

  rules {
    description = ["4"]
    rule_name   = "Video On Demand"
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
          full_url_cache = "off"
          ignore_case    = "off"
          scheme         = null
          query_string {
            action = null
            switch = "off"
            values = []
          }
        }
      }

      actions {
        name = "RangeOriginPull"
        range_origin_pull_parameters {
          switch = "on"
        }
      }

      sub_rules {
        description = ["4-1"]
        branches {
          condition = "$${http.request.file_extension} in ['php', 'jsp', 'asp', 'aspx']"
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
    }
  }

  rules {
    description = ["5"]
    rule_name   = "API Acceleration"
    branches {
      condition = "$${http.request.host} in ['aaa.makn.cn']"
      actions {
        name = "Cache"
        cache_parameters {
          no_cache {
            switch = "on"
          }
        }
      }

      actions {
        name = "SmartRouting"
        smart_routing_parameters {
          switch = "off"
        }
      }
    }
  }

  rules {
    description = ["6"]
    rule_name   = "WordPress Site"
    branches {
      condition = "$${http.request.host} in ['aaa.makn.cn']"
      sub_rules {
        description = ["6-1"]
        branches {
          condition = "$${http.request.file_extension} in ['gif', 'png', 'bmp', 'jpeg', 'tif', 'tiff', 'zip', 'exe', 'wmv', 'swf', 'mp3', 'wma', 'rar', 'css', 'flv', 'mp4', 'txt', 'ico', 'js']"
          actions {
            name = "Cache"
            cache_parameters {
              custom_time {
                cache_time           = 604800
                ignore_cache_control = "off"
                switch               = "on"
              }
            }
          }
        }

        branches {
          condition = "$${http.request.uri.path} in ['/']"
          actions {
            name = "Cache"
            cache_parameters {
              no_cache {
                switch = "on"
              }
            }
          }
        }

        branches {
          condition = "$${http.request.file_extension} in ['aspx', 'jsp', 'php', 'asp', 'do', 'dwr', 'cgi', 'fcgi', 'action', 'ashx', 'axd']"
          actions {
            name = "Cache"
            cache_parameters {
              no_cache {
                switch = "on"
              }
            }
          }
        }

        branches {
          condition = "$${http.request.uri.path} in ['/wp-admin/']"
          actions {
            name = "Cache"
            cache_parameters {
              no_cache {
                switch = "on"
              }
            }
          }
        }

        branches {
          condition = "*"
          actions {
            name = "Cache"
            cache_parameters {
              follow_origin {
                default_cache          = "on"
                default_cache_strategy = "on"
                default_cache_time     = 0
                switch                 = "on"
              }
            }
          }
        }
      }
    }
  }
}
```

## Using Filters to Filter Query Results

You can use the `filters` parameter to filter the rules returned from the API.

```hcl
resource "tencentcloud_teo_l7_acc_rule" "filtered" {
  zone_id = "zone-36bjhygh1bxe"

  # Filter rules by name
  filters {
    name   = "RuleName"
    values = ["Web Acceleration", "Live Video Streaming"]
  }

  rules {
    # ... rule configuration ...
  }
}
```

```hcl
resource "tencentcloud_teo_l7_acc_rule" "filtered_by_status" {
  zone_id = "zone-36bjhygh1bxe"

  # Filter rules by status
  filters {
    name   = "Status"
    values = ["Enabled"]
  }

  rules {
    # ... rule configuration ...
  }
}
```

```hcl
resource "tencentcloud_teo_l7_acc_rule" "multi_filter" {
  zone_id = "zone-36bjhygh1bxe"

  # Multiple filters (AND logic between different filters)
  filters {
    name   = "RuleName"
    values = ["Web Acceleration"]
  }
  filters {
    name   = "Status"
    values = ["Enabled"]
  }

  rules {
    # ... rule configuration ...
  }
}
```

**Note:** Multiple filter blocks with different names are combined with AND logic (all conditions must be met). Multiple values within a single filter are combined with OR logic (any value can match).

Supported filter names include:
- `rule-id`: Filter by rule ID
- `RuleName`: Filter by rule name
- `Status`: Filter by rule status
- Other fields supported by the DescribeL7AccRules API

Import

TEO l7 acc rule can be imported using the zone_id, e.g.

````
terraform import tencentcloud_teo_l7_acc_rule.example zone-36bjhygh1bxe
````
