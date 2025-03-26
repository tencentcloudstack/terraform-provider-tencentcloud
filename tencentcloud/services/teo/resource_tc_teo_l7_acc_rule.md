Provides a resource to create a teo l7_acc_rule

Example Usage

```hcl
resource "tencentcloud_teo_l7_acc_rule" "teo_l7_acc_rule" {
  zone_id = "zone-36bjhygh1bxe"
  rules {
    description = ["1"]
    rule_name   = "网站加速"
    status      = "disable"
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
    rule_name   = "音视频直播"
    status      = "enable"
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
    rule_name   = "大文件下载"
    status      = "enable"
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
    rule_name   = "音视频点播"
    status      = "enable"
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
    rule_name   = "API 加速"
    status      = "enable"
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
    rule_name   = "WordPress 建站"
    status      = "enable"
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
Import

teo l7_acc_rule can be imported using the zone_id, e.g.
````
terraform import tencentcloud_teo_l7_acc_rule.teo_l7_acc_rule zone-297z8rf93cfw
````