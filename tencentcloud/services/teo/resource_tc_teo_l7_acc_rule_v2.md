Provides a resource to create a TEO l7 acc rule v2

~> **NOTE:** Compared to tencentcloud_teo_l7_acc_rule, tencentcloud_teo_l7_acc_rule_v2 is simpler to use but is limited to managing a single rule and lacks the ability to maintain rule ordering. It is best suited for scenarios where you need to manage multiple rules independently and priority/sequencing is not a concern.

Example Usage

```hcl
resource "tencentcloud_teo_l7_acc_rule_v2" "example" {
  zone_id     = "zone-3fkff38fyw8s"
  description = ["description"]
  rule_name   = "Web Acceleration"
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

    actions {
      name = "Vary"
      vary_parameters {
        switch = "on"
      }
    }

    actions {
      name = "OriginAuthentication"
      origin_authentication_parameters {
        request_properties {
          type  = "Header"
          name  = "Authorization"
          value = "Bearer token123"
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

Using AdvancedOriginRouting, Shield, and SiteFailover actions

```hcl
resource "tencentcloud_teo_l7_acc_rule_v2" "example_failover" {
  zone_id     = "zone-3fkff38fyw8s"
  description = ["description"]
  rule_name   = "Web Acceleration Failover"
  status      = "enable"
  branches {
    condition = "$${http.request.host} in ['www.example.com']"
    actions {
      name = "AdvancedOriginRouting"
      advanced_origin_routing_parameters {
        direction = "MainlandChinaAndGlobalAdaptive"
      }
    }

    actions {
      name = "Shield"
      shield_parameters {
        shield_space_id = "shield-space-abc123"
      }
    }

    actions {
      name = "SiteFailover"
      site_failover_parameters {
        site_failover_status_codes = [500]
        site_failover_params {
          mode           = "FailoverToHost"
          origin         = "backup.example.com"
          origin_protocol = "https"
          https_origin_port = 443
          status_code    = 302
        }
      }
    }
  }
}
```

Using custom_variables and custom_variable_operations

```hcl
resource "tencentcloud_teo_l7_acc_rule_v2" "example_custom_vars" {
  zone_id     = "zone-3fkff38fyw8s"
  description  = ["description"]
  rule_name   = "Web Acceleration Custom Vars"
  status       = "enable"

  custom_variables {
    name          = "user.rule.var1"
    initial_value = "default_value"
    description   = "custom variable 1"
  }

  custom_variable_operations {
    branches {
      condition = "$${http.request.host} in ['www.example.com']"
      actions {
        name = "Set"
        set_parameters {
          name  = "user.rule.var1"
          value = "$${http.request.uri.path}"
        }
      }
      sub_rules {
        branches {
          condition = "lower($${http.request.file_extension}) in ['html']"
          actions {
            name = "Set"
            set_parameters {
              name  = "user.rule.var1"
              value = "static"
            }
          }
        }
        description = ["sub-rule-desc"]
      }
    }
    description = ["main-desc"]
  }

  branches {
    condition = "$${http.request.host} in ['www.example.com']"
  }
}
```

Using CustomAction with custom_action_parameters

```hcl
resource "tencentcloud_teo_l7_acc_rule_v2" "example_custom_action" {
  zone_id     = "zone-3fkff38fyw8s"
  description  = ["description"]
  rule_name   = "Web Acceleration Custom Action"
  status       = "enable"

  branches {
    condition = "$${http.request.host} in ['www.example.com']"
    actions {
      name = "CustomAction"
      custom_action_parameters {
        custom_actions {
          action_id = "action-001"
          parameters {
            name        = "param1"
            value_type  = "String"
            string_value = "value1"
          }
          parameters {
            name           = "param2"
            value_type     = "Integer"
            integer_value  = 42
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
